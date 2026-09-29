package cron

import (
	"context"
	"fmt"
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/gkarman/demo/internal/worker/cron")

// slogCronLogger adapts slog.Logger to the cron.Logger interface.
type slogCronLogger struct{ log *slog.Logger }

func (l *slogCronLogger) Info(msg string, keysAndValues ...any) {
	l.log.Info("cron: "+msg, keysAndValues...)
}

func (l *slogCronLogger) Error(err error, msg string, keysAndValues ...any) {
	l.log.Error("cron: "+msg, append(keysAndValues, "error", err)...)
}

type noopCronLogger struct{}

func (noopCronLogger) Info(_ string, _ ...any)           {}
func (noopCronLogger) Error(_ error, _ string, _ ...any) {}

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
	job := func(name string, fn func(context.Context) error) cron.Job {
		run := func() { w.runJob(name, fn) }
		return cron.NewChain(cron.SkipIfStillRunning(noopCronLogger{})).Then(cron.FuncJob(run))
	}

	if _, err := w.cron.AddJob("0 3 * * *", job("refresh_all_bloggers", w.refreshAllBloggers)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", job("poll_video_generations", w.pollVideoGenerations)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", job("poll_broll_generations", w.pollBrollGenerations)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", job("poll_compositions", w.pollCompositions)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", job("retry_pending_broll_submissions", w.retryPendingBrollSubmissions)); err != nil {
		return err
	}
	if _, err := w.cron.AddJob("*/1 * * * *", job("trigger_pending_compositions", w.triggerPendingCompositions)); err != nil {
		return err
	}
	return nil
}

// runJob выполняет задачу, логирует ошибку, записывает метрики и трейс.
// Каждый запуск — отдельный трейс: SQL и вызовы внешних API задачи станут его дочерними спанами.
func (w *Worker) runJob(name string, fn func(context.Context) error) {
	ctx, span := tracer.Start(w.ctx, "cron "+name, trace.WithNewRoot())
	defer span.End()

	start := time.Now()
	result := "error" // если fn запаникует, запуск засчитается как ошибка (панику перехватит cron.Recover)
	defer func() {
		jobRunsTotal.WithLabelValues(name, result).Inc()
		jobDuration.WithLabelValues(name).Observe(time.Since(start).Seconds())
	}()

	if err := fn(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.log.Error("cron job failed", "job", name, "error", err)
		return
	}

	result = "ok"
	jobLastSuccess.WithLabelValues(name).SetToCurrentTime()
}

func (w *Worker) refreshAllBloggers(ctx context.Context) error {
	w.log.Info("cron: refreshAllBloggers started")
	bloggerRepo := blogger.NewPostgres(w.db)
	videoSearcher := apifysearcher.NewVideoSearcher(w.apifyClient)
	fetchVideoCmd := command.NewFetchBloggerVideos(bloggerRepo, videoSearcher)

	refreshCmd := command.NewRefreshAllBloggers(bloggerRepo, fetchVideoCmd)

	if err := refreshCmd.Execute(ctx); err != nil {
		return fmt.Errorf("refresh all bloggers: %w", err)
	}
	return nil
}

func (w *Worker) pollBrollGenerations(ctx context.Context) error {
	w.log.Info("polling broll generations...")
	repo := blogger.NewPostgres(w.db)
	composeCmd := command.NewComposeFinalVideo(repo, w.videoComposer)
	pollCmd := command.NewPollBrollGenerations(repo, w.brollVideoGenerator, composeCmd)

	if err := pollCmd.Execute(ctx); err != nil {
		return fmt.Errorf("poll broll generations: %w", err)
	}
	w.log.Info("polling broll generations done")
	return nil
}

func (w *Worker) pollCompositions(ctx context.Context) error {
	w.log.Info("polling compositions...")
	repo := blogger.NewPostgres(w.db)

	disp := dispatcher.New()
	disp.Register(&bloggerDomain.VideoCompositionDone{}, bloggerHandlers.VideoCompositionDoneToRabbitHandler(w.publisher, w.log))

	pollCmd := command.NewPollCompositions(repo, w.videoComposer, disp)

	if err := pollCmd.Execute(ctx); err != nil {
		return fmt.Errorf("poll compositions: %w", err)
	}
	w.log.Info("polling compositions done")
	return nil
}

// Ошибки по отдельным видео только логируются: задача в целом отработала.
func (w *Worker) retryPendingBrollSubmissions(ctx context.Context) error {
	w.log.Info("cron: retryPendingBrollSubmissions started")
	repo := blogger.NewPostgres(w.db)
	videoIDs, err := repo.ListVideosWithPendingBrollSegments(ctx)
	if err != nil {
		return fmt.Errorf("list videos with pending broll segments: %w", err)
	}
	if len(videoIDs) == 0 {
		return nil
	}
	w.log.Info("retrying pending broll submissions", "videos", len(videoIDs))
	submitCmd := command.NewSubmitBrollGenerations(repo, w.brollVideoGenerator)
	for _, videoID := range videoIDs {
		if err := submitCmd.Run(ctx, reqdto.SubmitBrollGenerations{VideoID: videoID}); err != nil {
			w.log.Error("failed to submit broll generations", "video_id", videoID, "error", err)
		}
	}
	return nil
}

// Ошибки по отдельным видео только логируются: задача в целом отработала.
func (w *Worker) triggerPendingCompositions(ctx context.Context) error {
	w.log.Info("cron: triggerPendingCompositions started")
	repo := blogger.NewPostgres(w.db)
	videoIDs, err := repo.ListVideosReadyToCompose(ctx)
	if err != nil {
		return fmt.Errorf("list videos ready to compose: %w", err)
	}
	if len(videoIDs) == 0 {
		return nil
	}
	w.log.Info("found videos ready to compose", "count", len(videoIDs))
	composeCmd := command.NewComposeFinalVideo(repo, w.videoComposer)
	for _, videoID := range videoIDs {
		if err := composeCmd.Run(ctx, reqdto.ComposeFinalVideo{VideoID: videoID}); err != nil {
			w.log.Error("failed to compose video", "video_id", videoID, "error", err)
		} else {
			w.log.Info("composition triggered", "video_id", videoID)
		}
	}
	return nil
}

func (w *Worker) pollVideoGenerations(ctx context.Context) error {
	w.log.Info("polling video generations...")
	bloggerRepo := blogger.NewPostgres(w.db)

	disp := dispatcher.New()
	disp.Register(&bloggerDomain.VideoGenerationDone{}, bloggerHandlers.VideoGenerationDoneToRabbitHandler(w.publisher, w.log))
	disp.Register(&bloggerDomain.VideoGenerationError{}, bloggerHandlers.VideoGenerationErrorToRabbitHandler(w.publisher, w.log))

	composeCmd := command.NewComposeFinalVideo(bloggerRepo, w.videoComposer)
	pollCmd := command.NewPollVideoGenerations(bloggerRepo, w.videoGenerator, w.storage, disp, composeCmd)

	if err := pollCmd.Execute(ctx); err != nil {
		return fmt.Errorf("poll video generations: %w", err)
	}
	w.log.Info("polling video generations done")
	return nil
}
