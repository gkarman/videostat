package notify

import (
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/worker/notify/handlers"
)

func NewRouterWithHandlers(log *slog.Logger) *Router {
	r := NewRouter(log)

	carCreated := handlers.NewCarCreatedHandler(log)
	carUpdated := handlers.NewCarUpdatedHandler(log)

	r.Register(events.EventCarCreatedV1, carCreated.Handle)
	r.Register(events.EventCarUpdatedV1, carUpdated.Handle)

	return r
}