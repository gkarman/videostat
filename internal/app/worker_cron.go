package app

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/platform"
	cronworker "github.com/gkarman/demo/internal/worker/cron"
)

func NewWorkerCron(ctx context.Context) (*cronworker.Worker, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := platform.NewLogger(cfg)

	log.Info("db connect...")
	db, err := platform.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	log.Info("db connected")

	worker, err := cronworker.New(log, db)
	if err != nil {
		return nil, fmt.Errorf("create worker cron: %w", err)
	}

	return worker, nil
}
