package platform

import (
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/mq"
)

func NewRabbitConsumer(cfg *config.Config, log *slog.Logger) (*mq.RabbitConsumer, error) {
	configRabbit := mq.Config{
		User:           cfg.RabbitMQ.User,
		Password:       cfg.RabbitMQ.Password,
		Host:           cfg.RabbitMQ.Host,
		Port:           cfg.RabbitMQ.Port,
		Exchange:       cfg.RabbitMQ.Exchange,
		ReconnectDelay: time.Duration(cfg.RabbitMQ.ReconnectDelay) * time.Second,
	}

	consumer := mq.NewRabbitConsumer(
		configRabbit,
		"worker.notify",
		[]string{
			"car.#",
			"user.#",
		},
		log,
	)

	return consumer, nil
}
