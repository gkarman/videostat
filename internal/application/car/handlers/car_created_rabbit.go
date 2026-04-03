package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/events/mappers"
)

func CarCreatedToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*car.Created)
		if !ok {
			log.Error("invalid event type for car.created (rabbit)")
			return
		}

		msg := mappers.MapCarCreated(event)
		body, err := json.Marshal(msg)
		if err != nil {
			log.Error("marshal failed in CarCreatedToRabbitHandler", "err", err)
			return
		}

		err = publisher.Publish(ctx, events.EventCarCreatedV1, body)
		if err != nil {
			log.Error("failed to publish to rabbitmq", "err", err)
			return
		}

		log.Info("event published to rabbitmq", "event", msg.EventID)
	}
}
