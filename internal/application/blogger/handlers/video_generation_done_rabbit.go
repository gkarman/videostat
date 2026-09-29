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

func VideoGenerationDoneToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.VideoGenerationDone)
		if !ok {
			log.Error("invalid event type for blogger.VideoGenerationDone", "event", e)
			return
		}

		body, err := json.Marshal(mappers.MapVideoGenerationDone(event))
		if err != nil {
			log.Error("marshal VideoGenerationDone failed", "error", err)
			return
		}

		if err = publisher.Publish(ctx, events.EventVideoGenerationDoneV1, body); err != nil {
			log.Error("publish VideoGenerationDone failed", "error", err)
		}
	}
}

func VideoCompositionDoneToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.VideoCompositionDone)
		if !ok {
			log.Error("invalid event type for blogger.VideoCompositionDone", "event", e)
			return
		}

		body, err := json.Marshal(mappers.MapVideoCompositionDone(event))
		if err != nil {
			log.Error("marshal VideoCompositionDone failed", "error", err)
			return
		}

		if err = publisher.Publish(ctx, events.EventVideoCompositionDoneV1, body); err != nil {
			log.Error("publish VideoCompositionDone failed", "error", err)
		}
	}
}

func VideoGenerationErrorToRabbitHandler(publisher application.Publisher, log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*blogger.VideoGenerationError)
		if !ok {
			log.Error("invalid event type for blogger.VideoGenerationError", "event", e)
			return
		}

		body, err := json.Marshal(mappers.MapVideoGenerationError(event))
		if err != nil {
			log.Error("marshal VideoGenerationError failed", "error", err)
			return
		}

		if err = publisher.Publish(ctx, events.EventVideoGenerationErrorV1, body); err != nil {
			log.Error("publish VideoGenerationError failed", "error", err)
		}
	}
}
