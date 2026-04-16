package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	log  *slog.Logger
	db   *pgxpool.Pool
	cron *cron.Cron
	ctx  context.Context
	apifyClient *apify.Client
}

func New(log *slog.Logger, db *pgxpool.Pool, apifyClient *apify.Client) (*Worker, error) {
	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &Worker{
		log:  log,
		db:   db,
		cron: c,
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
	_, err := w.cron.AddFunc("0 3 * * *", w.refreshAllBloggers)

	return err
}

func (w *Worker) refreshAllBloggers() {
	bloggerRepo := blogger.NewPostgres(w.db)
	videoSearcher := apify.NewVideoSearcher(w.apifyClient)
	fetchVideoCmd := command.NewFetchBloggerVideos(bloggerRepo, videoSearcher)

	refreshCmd := command.NewRefreshAllBloggers(
		bloggerRepo,
		fetchVideoCmd,
	)

	err := refreshCmd.Execute(w.ctx)
	if err != nil {
		w.log.Error("Failed to refresh all Bloggers.", err)
	}
}
