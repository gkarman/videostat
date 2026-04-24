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

func BloggerCreatedToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.Created)
		if !ok {
			log.Error("invalid event type for blogger.created (rabbit)")
			return
		}

		msg := mappers.MapBloggerCreated(event)
		body, err := json.Marshal(msg)
		if err != nil {
			log.Error("marshal failed in BloggerCreatedToRabbitHandler", "err", err)
			return
		}

		err = publisher.Publish(ctx, events.EventBloggerCreatedV1, body)
		if err != nil {
			log.Error("failed to publish to rabbitmq", "err", err)
			return
		}

		log.Info("event published to rabbitmq", "event", msg.EventID)
	}
}
