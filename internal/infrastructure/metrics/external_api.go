package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Внешние API (генерация видео, LLM) отвечают от долей секунды до минуты.
var externalAPIDurationBuckets = []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20, 30, 60, 120}

var (
	externalAPIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_api_requests_total",
			Help: "Количество запросов во внешние API. status — HTTP-код ответа или error, если ответа не было (таймаут, сеть).",
		},
		[]string{"provider", "status"},
	)

	externalAPIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_request_duration_seconds",
			Help:    "Время до получения ответа от внешнего API в секундах.",
			Buckets: externalAPIDurationBuckets,
		},
		[]string{"provider"},
	)
)

// instrumentedTransport — обёртка над http.RoundTripper: каждый запрос через http.Client
// проходит через RoundTrip, поэтому здесь можно замерить все вызовы провайдера в одном месте.
type instrumentedTransport struct {
	provider string
	next     http.RoundTripper
}

// InstrumentedTransport возвращает транспорт для http.Client, который пишет метрики
// с меткой provider. Использование: &http.Client{Transport: metrics.InstrumentedTransport("heygen")}.
func InstrumentedTransport(provider string) http.RoundTripper {
	return &instrumentedTransport{provider: provider, next: http.DefaultTransport}
}

func (t *instrumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)

	status := "error"
	if err == nil {
		status = strconv.Itoa(resp.StatusCode)
	}
	externalAPIRequestsTotal.WithLabelValues(t.provider, status).Inc()
	externalAPIRequestDuration.WithLabelValues(t.provider).Observe(time.Since(start).Seconds())

	return resp, err
}
