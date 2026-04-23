package notify

import (
	"log/slog"

	"github.com/gkarman/demo/internal/worker"
)

func NewRouterWithHandlers(log *slog.Logger) *worker.Router {
	r := worker.NewRouter(log)
	return r
}