package app

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/platform"
	"github.com/gkarman/demo/internal/worker/core"
)

func NewWorkerCore(ctx context.Context) (*core.Worker, error) {
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
	consumer, err := platform.NewRabbitConsumer(
		cfg,
		log,
		"worker.core",
		[]string{
			"blogger.#",
		})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("init rabbit consumer: %w", err)
	}

	log.Info("rabbit consumer connected")

	log.Info("rabbit publisher connect...")
	publisher, err := platform.NewRabbitPublisher(cfg)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("init rabbit publisher: %w", err)
	}
	log.Info("rabbit publisher connected")

	apifyClient := platform.NewApifyClient(cfg)
	analyzer := platform.NewAssemblyAIAnalyzer(cfg)
	promptGenerator := platform.NewClaudePromptGenerator(cfg)
	router := core.NewRouterWithHandlers(log, db, apifyClient, analyzer, promptGenerator, publisher)
	worker := core.New(log, consumer, router)

	return worker, nil
}
