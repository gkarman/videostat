package worker

import (
	"context"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type Handler func(context.Context, []byte) error

type Router struct {
	log      *slog.Logger
	handlers map[string]Handler
}

func NewRouter(log *slog.Logger) *Router {
	return &Router{
		log:      log,
		handlers: make(map[string]Handler),
	}
}

func (r *Router) Register(eventType string, handler Handler) {
	r.handlers[eventType] = handler
}

func (r *Router) Handle(eventType string, body []byte) error {
	h, ok := r.handlers[eventType]
	if !ok {
		return nil
	}

	ctx := logger.WithLogger(context.Background(), r.log)
	return h(ctx, body)
}
