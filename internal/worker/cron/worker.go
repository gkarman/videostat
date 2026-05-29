package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	bloggerHandlers "github.com/gkarman/demo/internal/application/blogger/handlers"
	bloggerDomain "github.com/gkarman/demo/internal/domain/blogger"
	sharedapify "github.com/gkarman/demo/internal/infrastructure/apify"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	apifysearcher "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

// slogCronLogger adapts slog.Logger to the cron.Logger interface.
type slogCronLogger struct{ log *slog.Logger }

func (l *slogCronLogger) Info(msg string, keysAndValues ...any) {
	l.log.Info("cron: "+msg, keysAndValues...)
}

func (l *slogCronLogger) Error(err error, msg string, keysAndValues ...any) {
	l.log.Error("cron: "+msg, append(keysAndValues, "error", err)...)
}

type noopCronLogger struct{}

func (noopCronLogger) Info(_ string, _ ...any)             {}
func (noopCronLogger) Error(_ error, _ string, _ ...any)   {}

type Worker struct {
	log                 *slog.Logger
	db                  *pgxpool.Pool
	cron                *cron.Cron
	ctx                 context.Context
	apifyClient         *sharedapify.Client
	videoGenerator      application.VideoGenerator
	brollVideoGenerator application.BrollVideoGenerator
	videoComposer       application.VideoComposer
	storage             application.Storage
	publisher           application.Publisher
}

func New(
	log *slog.Logger,
	db *pgxpool.Pool,
	apifyClient *sharedapify.Client,
	videoGenerator application.VideoGenerator,
	brollVideoGenerator application.BrollVideoGenerator,
	videoComposer application.VideoComposer,
	storage application.Storage,
	publisher application.Publisher,
) (*Worker, error) {
	cronLog := &slogCronLogger{log: log}
	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithChain(cron.Recover(cronLog)),
	)

	return &Worker{
		log:                 log,
		db:                  db,
		cron:                c,
		apifyClient:         apifyClient,
		videoGenerator:      videoGenerator,
		brollVideoGenerator: brollVideoGenerator,
		videoComposer:       videoComposer,
		storage:             storage,
		publisher:           publisher,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	w.ctx = logger.WithLogger(ctx, w.log)

	if err := w.registerJobs(); err != nil {
		return err
	}

	w.cron.Start()
	<-ctx.Done()

	w.cron.Stop()
	return nil
}

func (w *Worker) registerJobs() error {
	skip := func(fn func()) cron.Job {
		return cron.NewChain(cron.SkipIfStillRunning(noopCronLogger{})).Then(cron.FuncJob(fn))
	}

	if _, err := w.cron.AddJob("0 3 * * *", skip(w.refreshAllBloggers)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", skip(w.pollVideoGenerations)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", skip(w.pollBrollGenerations)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", skip(w.pollCompositions)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", skip(w.retryPendingBrollSubmissions)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", skip(w.triggerPendingCompositions)); err != nil {
		return err
	}
	return nil
}

func (w *Worker) refreshAllBloggers() {
	w.log.Info("cron: refreshAllBloggers started")
	bloggerRepo := blogger.NewPostgres(w.db)
	videoSearcher := apifysearcher.NewVideoSearcher(w.apifyClient)
	fetchVideoCmd := command.NewFetchBloggerVideos(bloggerRepo, videoSearcher)

	refreshCmd := command.NewRefreshAllBloggers(bloggerRepo, fetchVideoCmd)

	if err := refreshCmd.Execute(w.ctx); err != nil {
		w.log.Error("Failed to refresh all Bloggers.", err)
	}
}

func (w *Worker) pollBrollGenerations() {
	w.log.Info("polling broll generations...")
	repo := blogger.NewPostgres(w.db)
	composeCmd := command.NewComposeFinalVideo(repo, w.videoComposer)
	pollCmd := command.NewPollBrollGenerations(repo, w.brollVideoGenerator, composeCmd)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("failed to poll broll generations", "error", err)
	} else {
		w.log.Info("polling broll generations done")
	}
}

func (w *Worker) pollCompositions() {
	w.log.Info("polling compositions...")
	repo := blogger.NewPostgres(w.db)

	disp := dispatcher.New()
	disp.Register(&bloggerDomain.VideoCompositionDone{}, bloggerHandlers.VideoCompositionDoneToRabbitHandler(w.publisher, w.log))

	pollCmd := command.NewPollCompositions(repo, w.videoComposer, disp)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("failed to poll compositions", "error", err)
	} else {
		w.log.Info("polling compositions done")
	}
}

func (w *Worker) retryPendingBrollSubmissions() {
	w.log.Info("cron: retryPendingBrollSubmissions started")
	repo := blogger.NewPostgres(w.db)
	videoIDs, err := repo.ListVideosWithPendingBrollSegments(w.ctx)
	if err != nil {
		w.log.Error("failed to list videos with pending broll segments", "error", err)
		return
	}
	if len(videoIDs) == 0 {
		return
	}
	w.log.Info("retrying pending broll submissions", "videos", len(videoIDs))
	submitCmd := command.NewSubmitBrollGenerations(repo, w.brollVideoGenerator)
	for _, videoID := range videoIDs {
		if err := submitCmd.Run(w.ctx, reqdto.SubmitBrollGenerations{VideoID: videoID}); err != nil {
			w.log.Error("failed to submit broll generations", "videoID", videoID, "error", err)
		}
	}
}

func (w *Worker) triggerPendingCompositions() {
	w.log.Info("cron: triggerPendingCompositions started")
	repo := blogger.NewPostgres(w.db)
	videoIDs, err := repo.ListVideosReadyToCompose(w.ctx)
	if err != nil {
		w.log.Error("failed to list videos ready to compose", "error", err)
		return
	}
	if len(videoIDs) == 0 {
		return
	}
	w.log.Info("found videos ready to compose", "count", len(videoIDs))
	composeCmd := command.NewComposeFinalVideo(repo, w.videoComposer)
	for _, videoID := range videoIDs {
		if err := composeCmd.Run(w.ctx, reqdto.ComposeFinalVideo{VideoID: videoID}); err != nil {
			w.log.Error("failed to compose video", "videoID", videoID, "error", err)
		} else {
			w.log.Info("composition triggered", "videoID", videoID)
		}
	}
}

func (w *Worker) pollVideoGenerations() {
	w.log.Info("polling video generations...")
	bloggerRepo := blogger.NewPostgres(w.db)

	disp := dispatcher.New()
	disp.Register(&bloggerDomain.VideoGenerationDone{}, bloggerHandlers.VideoGenerationDoneToRabbitHandler(w.publisher, w.log))
	disp.Register(&bloggerDomain.VideoGenerationError{}, bloggerHandlers.VideoGenerationErrorToRabbitHandler(w.publisher, w.log))

	composeCmd := command.NewComposeFinalVideo(bloggerRepo, w.videoComposer)
	pollCmd := command.NewPollVideoGenerations(bloggerRepo, w.videoGenerator, w.storage, disp, composeCmd)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("failed to poll video generations", "error", err)
	} else {
		w.log.Info("polling video generations done")
	}
}
