package notify

import (
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/notifier"
	"github.com/gkarman/demo/internal/worker"
	notify_handlers "github.com/gkarman/demo/internal/worker/notify/handlers"
)

func NewRouterWithHandlers(log *slog.Logger, n *notifier.TelegramNotifier) *worker.Router {
	r := worker.NewRouter(log)

	r.Register(events.EventVideoGenerationDoneV1, notify_handlers.NewVideoGenerationDoneHandler(n, log).Handle)
	r.Register(events.EventVideoGenerationErrorV1, notify_handlers.NewVideoGenerationErrorHandler(n, log).Handle)

	return r
}
