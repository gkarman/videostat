package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	sharedapify "github.com/gkarman/demo/internal/infrastructure/apify"
	apifysearcher "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	log            *slog.Logger
	db             *pgxpool.Pool
	cron           *cron.Cron
	ctx            context.Context
	apifyClient    *sharedapify.Client
	videoGenerator application.VideoGenerator
	storage        application.Storage
}

func New(
	log *slog.Logger,
	db *pgxpool.Pool,
	apifyClient *sharedapify.Client,
	videoGenerator application.VideoGenerator,
	storage application.Storage,
) (*Worker, error) {
	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &Worker{
		log:            log,
		db:             db,
		cron:           c,
		apifyClient:    apifyClient,
		videoGenerator: videoGenerator,
		storage:        storage,
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

func (w *Worker) pollVideoGenerations() {
	bloggerRepo := blogger.NewPostgres(w.db)
	pollCmd := command.NewPollVideoGenerations(bloggerRepo, w.videoGenerator, w.storage)

	if err := pollCmd.Execute(w.ctx); err != nil {
		w.log.Error("Failed to poll video generations", "error", err)
	}
}
