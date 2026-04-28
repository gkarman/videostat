package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/events/mappers"
)

func VideoSourceFoundToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.VideoSourceFound)
		if !ok {
			log.Error("invalid event type for blogger.VideoSourceFound (rabbit)", "event", e)
			return
		}

		integrationEvent := mappers.MapVideoSourceFound(event)
		body, err := json.Marshal(integrationEvent)
		if err != nil {
			log.Error("marshal failed in VideoSourceFoundToRabbitHandler", "err", err)
			return
		}

		err = publisher.Publish(ctx, events.EventVideoSourceFoundV1, body)
		if err != nil {
			log.Error("failed to publish to rabbitmq", "err", err)
			return
		}

		log.Debug("event published to rabbitmq", "event", integrationEvent)
	}
}
