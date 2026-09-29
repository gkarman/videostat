package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Имена полей логов — одинаковые во всех сервисах и в snake_case (как в событиях RabbitMQ).
// Иначе запрос в Loki {env="local"} | json | video_id="..." найдёт не все строки.
// Формат ключей проверяет линтер sloglint (.golangci.yml).
const (
	KeyVideoID   = "video_id"
	KeyBloggerID = "blogger_id"
	KeyEventType = "event_type"
	KeyEventID   = "event_id"
	KeyTraceID   = "trace_id"
)

type fieldKey string

// WithField добавляет поле в логгер контекста: все логи, которые дальше по цепочке вызовов
// берут логгер через FromContext(ctx), получат это поле автоматически.
// Если то же поле с тем же значением уже добавлено выше — повторно не добавляет,
// чтобы в JSON не было двух одинаковых ключей.
func WithField(ctx context.Context, key, value string) context.Context {
	if value == "" {
		return ctx
	}
	if existing, ok := ctx.Value(fieldKey(key)).(string); ok && existing == value {
		return ctx
	}

	ctx = context.WithValue(ctx, fieldKey(key), value)
	return WithLogger(ctx, FromContext(ctx).With(key, value))
}

// WithTraceID добавляет в логгер контекста trace_id текущего спана. Вызывается там, где начинается
// спан (HTTP-запрос, обработка события, запуск cron, сообщение в боте): все логи внутри получат
// trace_id, и в Grafana из строки лога можно перейти в трейс, а из трейса — к его логам.
func WithTraceID(ctx context.Context) context.Context {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ctx // трейсинг выключен или спана нет
	}
	return WithField(ctx, KeyTraceID, sc.TraceID().String())
}
