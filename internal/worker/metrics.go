package worker

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Возможные значения метки result.
const (
	resultOK      = "ok"
	resultError   = "error"
	resultSkipped = "skipped" // для события нет обработчика
)

// eventTypeUnhandled — общее значение метки для событий без обработчика:
// тип приходит извне, и писать его как есть — риск раздуть число временных рядов.
const eventTypeUnhandled = "unhandled"

// Обработка сообщения может ходить во внешние API (анализ, генерация видео),
// поэтому корзины — от 10 мс до 5 минут, а не стандартные до 10 с.
var messageDurationBuckets = []float64{0.01, 0.05, 0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120, 300}

var (
	messagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_messages_total",
			Help: "Количество обработанных сообщений из очереди.",
		},
		[]string{"event_type", "result"},
	)

	messageDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_message_duration_seconds",
			Help:    "Длительность обработки сообщения из очереди в секундах.",
			Buckets: messageDurationBuckets,
		},
		[]string{"event_type"},
	)
)
