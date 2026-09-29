package worker

import (
	"context"
	"encoding/json"
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
	ctx := r.eventContext(eventType, body)
	err := h(ctx, body)

	result := resultOK
	if err != nil {
		result = resultError
	}
	messagesTotal.WithLabelValues(eventType, result).Inc()
	messageDuration.WithLabelValues(eventType).Observe(time.Since(start).Seconds())

	return err
}

// eventContext кладёт в логгер поля события (тип, id, video_id / blogger_id).
// Все логи обработчика и вызванных им команд получат их автоматически — по video_id
// в Loki находится весь путь видео через api, worker_core, worker_cron и worker_notify.
func (r *Router) eventContext(eventType string, body []byte) context.Context {
	var ids struct {
		EventID   string `json:"event_id"`
		VideoID   string `json:"video_id"`
		BloggerID string `json:"blogger_id"`
	}
	// Ошибку разбора не обрабатываем: без этих полей обработчик всё равно отработает,
	// а невалидное тело он сам отклонит при своём разборе.
	_ = json.Unmarshal(body, &ids)

	ctx := logger.WithLogger(context.Background(), r.log)
	ctx = logger.WithField(ctx, logger.KeyEventType, eventType)
	ctx = logger.WithField(ctx, logger.KeyEventID, ids.EventID)
	ctx = logger.WithField(ctx, logger.KeyVideoID, ids.VideoID)
	ctx = logger.WithField(ctx, logger.KeyBloggerID, ids.BloggerID)
	return ctx
}
