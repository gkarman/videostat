package app

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/platform"
	"github.com/gkarman/demo/internal/worker/notify"
)

func NewWorkerNotify(ctx context.Context) (*notify.Worker, error) {
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

	log.Info("rabbit consumer connect...")
	consumer, err := platform.NewRabbitConsumer(cfg, log)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("init rabbit consumer: %w", err)
	}
	log.Info("rabbit consumer connected")

	router := notify.NewRouterWithHandlers(log)
	worker := notify.New(log, consumer, router)

	return worker, nil
}