package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

// ObserveWithTrace записывает значение в гистограмму и, если в ctx есть сохраняемый спан,
// прикладывает к нему exemplar — «пример» с trace_id. На графике в Grafana такие примеры
// видны точками: клик по точке на p95 открывает трейс именно этого медленного запроса.
func ObserveWithTrace(ctx context.Context, o prometheus.Observer, value float64) {
	sc := trace.SpanContextFromContext(ctx)
	if eo, ok := o.(prometheus.ExemplarObserver); ok && sc.IsSampled() {
		eo.ObserveWithExemplar(value, prometheus.Labels{"trace_id": sc.TraceID().String()})
		return
	}
	o.Observe(value)
}

// Handler — /metrics с поддержкой exemplars. Стандартный promhttp.Handler() отдаёт старый
// текстовый формат, в котором exemplars нет; они есть только в формате OpenMetrics,
// поэтому он включён явно (Prometheus сам запрашивает OpenMetrics, если его поддерживают).
func Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}),
	)
}
