package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	log  *slog.Logger
	db   *pgxpool.Pool
	cron *cron.Cron
	ctx  context.Context
}

func New(log *slog.Logger, db *pgxpool.Pool) (*Worker, error) {
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
	_, err := w.cron.AddFunc("@hourly", w.hourlyJob)
	return err
}

func (w *Worker) hourlyJob() {
	w.log.Info("hourly job started")
	// use-case
	w.log.Info("hourly job finished")
}
