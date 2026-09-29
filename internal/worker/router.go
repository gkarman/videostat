package worker

import (
	"context"
	"log/slog"
	"time"

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
		r.log.Debug("no handler for event type", "event_type", eventType)
		messagesTotal.WithLabelValues(eventTypeUnhandled, resultSkipped).Inc()
		return nil
	}

	start := time.Now()
	ctx := logger.WithLogger(context.Background(), r.log)
	err := h(ctx, body)

	result := resultOK
	if err != nil {
		result = resultError
	}
	messagesTotal.WithLabelValues(eventType, result).Inc()
	messageDuration.WithLabelValues(eventType).Observe(time.Since(start).Seconds())

	return err
}
