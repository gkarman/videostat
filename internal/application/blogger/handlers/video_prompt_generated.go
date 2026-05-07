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

func VideoPromptGeneratedToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.VideoPromptGenerated)
		if !ok {
			log.Error("invalid event type for blogger.VideoPromptGenerated (rabbit)", "event", e)
			return
		}

		integrationEvent := mappers.MapVideoPromptGenerated(event)
		body, err := json.Marshal(integrationEvent)
		if err != nil {
			log.Error("marshal failed in VideoPromptGeneratedToRabbitHandler", "err", err)
			return
		}

		if err = publisher.Publish(ctx, events.EventVideoPromptGeneratedV1, body); err != nil {
			log.Error("failed to publish to rabbitmq", "err", err)
			return
		}

		log.Debug("event published to rabbitmq", "event", integrationEvent)
	}
}
