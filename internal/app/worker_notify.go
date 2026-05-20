package app

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/notifier"
	"github.com/gkarman/demo/internal/platform"
	"github.com/gkarman/demo/internal/worker/notify"
)

func NewWorkerNotify(ctx context.Context) (*notify.Worker, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := platform.NewLogger(cfg)

	log.Info("rabbit consumer connect...")
	consumer, err := platform.NewRabbitConsumer(
		cfg,
		log,
		"worker.notify",
		[]string{
			"car.#",
			"user.#",
			"blogger.#",
		})
	if err != nil {
		return nil, fmt.Errorf("init rabbit consumer: %w", err)
	}
	log.Info("rabbit consumer connected")

	tgNotifier, err := notifier.NewTelegramNotifier(cfg.TelegramBot.Token)
	if err != nil {
		return nil, fmt.Errorf("init telegram notifier: %w", err)
	}

	router := notify.NewRouterWithHandlers(log, tgNotifier)
	worker := notify.New(log, consumer, router)

	return worker, nil
}
