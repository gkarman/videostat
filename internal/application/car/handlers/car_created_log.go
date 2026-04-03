package handlers

import (
	"context"
	"log/slog"

	"github.com/gkarman/demo/internal/domain/car"
)

func CarCreatedLogHandler(log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {
		event, ok := e.(*car.Created)
		if !ok {
			log.Error("invalid event type for car.created")
			return
		}

		log.Info("car created",
			"id", event.ID,
			"name", event.Name,
		)
	}
}