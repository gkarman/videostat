package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command"
	bloggerHandlers "github.com/gkarman/demo/internal/application/blogger/handlers"
	bloggerDomain "github.com/gkarman/demo/internal/domain/blogger"
	sharedapify "github.com/gkarman/demo/internal/infrastructure/apify"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	apifysearcher "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

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
	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
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
	w.ctx = ctx

	if err := w.registerJobs(); err != nil {
		return err
	}

	w.cron.Start()
	<-ctx.Done()

	w.cron.Stop()
	return nil
}

func (w *Worker) registerJobs() error {
	if _, err := w.cron.AddFunc("0 3 * * *", w.refreshAllBloggers); err != nil {
		return err
	}
	if _, err := w.cron.AddFunc("*/3 * * * *", w.pollVideoGenerations); err != nil {
		return err
	}
	if _, err := w.cron.AddFunc("*/1 * * * *", w.pollBrollGenerations); err != nil {
		return err
	}
	if _, err := w.cron.AddFunc("*/1 * * * *", w.pollCompositions); err != nil {
		return err
	}
	return nil
}

func (w *Worker) refreshAllBloggers() {
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
	pollCmd := command.NewPollCompositions(repo, w.videoComposer)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("failed to poll compositions", "error", err)
	} else {
		w.log.Info("polling compositions done")
	}
}

func (w *Worker) pollVideoGenerations() {
	bloggerRepo := blogger.NewPostgres(w.db)

	disp := dispatcher.New()
	disp.Register(&bloggerDomain.VideoGenerationDone{}, bloggerHandlers.VideoGenerationDoneToRabbitHandler(w.publisher, w.log))
	disp.Register(&bloggerDomain.VideoGenerationError{}, bloggerHandlers.VideoGenerationErrorToRabbitHandler(w.publisher, w.log))

	pollCmd := command.NewPollVideoGenerations(bloggerRepo, w.videoGenerator, w.storage, disp)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("Failed to poll video generations", "error", err)
	}
}
