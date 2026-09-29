package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/gkarman/demo/internal/worker")

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
	ids := parseEventIDs(body)
	ctx := r.eventContext(eventType, ids)

	// Спан на обработку сообщения. Пока каждое сообщение — новый трейс;
	// чтобы продолжать трейс отправителя, нужен traceparent в заголовках AMQP (шаг T2).
	ctx, span := tracer.Start(ctx, "handle "+eventType,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "rabbitmq"),
			attribute.String("event_id", ids.EventID),
			// Те же имена, что в логах: в Tempo ищется { span.video_id = "..." }
			attribute.String(logger.KeyVideoID, ids.VideoID),
			attribute.String(logger.KeyBloggerID, ids.BloggerID),
		),
	)
	defer span.End()

	err := h(ctx, body)

	result := resultOK
	if err != nil {
		result = resultError
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	messagesTotal.WithLabelValues(eventType, result).Inc()
	messageDuration.WithLabelValues(eventType).Observe(time.Since(start).Seconds())

	return err
}

// eventIDs — идентификаторы, которые есть в теле событий (см. contracts/events).
type eventIDs struct {
	EventID   string `json:"event_id"`
	VideoID   string `json:"video_id"`
	BloggerID string `json:"blogger_id"`
}

func parseEventIDs(body []byte) eventIDs {
	var ids eventIDs
	// Ошибку разбора не обрабатываем: без этих полей обработчик всё равно отработает,
	// а невалидное тело он сам отклонит при своём разборе.
	_ = json.Unmarshal(body, &ids)
	return ids
}

// eventContext кладёт в логгер поля события (тип, id, video_id / blogger_id).
// Все логи обработчика и вызванных им команд получат их автоматически — по video_id
// в Loki находится весь путь видео через api, worker_core, worker_cron и worker_notify.
func (r *Router) eventContext(eventType string, ids eventIDs) context.Context {
	ctx := logger.WithLogger(context.Background(), r.log)
	ctx = logger.WithField(ctx, logger.KeyEventType, eventType)
	ctx = logger.WithField(ctx, logger.KeyEventID, ids.EventID)
	ctx = logger.WithField(ctx, logger.KeyVideoID, ids.VideoID)
	ctx = logger.WithField(ctx, logger.KeyBloggerID, ids.BloggerID)
	return ctx
}
