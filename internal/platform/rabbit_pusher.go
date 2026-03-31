package platform

import (
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/mq"
)

func NewRabbitPublisher(cfg *config.Config) (*mq.RabbitPublisher, error) {
	configRabbit := mq.Config{
		User:           cfg.RabbitMQ.User,
		Password:       cfg.RabbitMQ.Password,
		Host:           cfg.RabbitMQ.Host,
		Port:           cfg.RabbitMQ.Port,
		Exchange:       cfg.RabbitMQ.Exchange,
		ReconnectDelay: time.Duration(cfg.RabbitMQ.ReconnectDelay) * time.Second,
	}

	publisher, err := mq.NewRabbitPublisher(configRabbit)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
