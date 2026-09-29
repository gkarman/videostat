package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/gkarman/demo/internal/infrastructure/mq")

// headersCarrier позволяет OpenTelemetry писать и читать traceparent в заголовках AMQP-сообщения.
// Для HTTP то же самое делает otelhttp с HTTP-заголовками, для RabbitMQ готовой обёртки нет —
// поэтому адаптер свой (реализует propagation.TextMapCarrier).
type headersCarrier amqp.Table

func (c headersCarrier) Get(key string) string {
	v, _ := c[key].(string)
	return v
}

func (c headersCarrier) Set(key, value string) {
	c[key] = value
}

func (c headersCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
