package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gkarman/demo/internal/infrastructure/metrics"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// promauto сразу регистрирует метрики в стандартном реестре,
// который отдаёт metrics.Handler() на /metrics.
var (
	// Counter: только растёт. Метки позволяют разрезать по методу, маршруту и статусу.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество обработанных HTTP-запросов.",
		},
		[]string{"method", "route", "status"},
	)

	// Histogram: раскладывает длительность запросов по корзинам (buckets).
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Длительность обработки HTTP-запроса в секундах.",
			Buckets: prometheus.DefBuckets, // 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
		},
		[]string{"method", "route"},
	)

	// Gauge: текущее значение, растёт и падает.
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Количество HTTP-запросов, которые обрабатываются прямо сейчас.",
		},
	)
)

func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpRequestsInFlight.Inc()
			defer httpRequestsInFlight.Dec()

			start := time.Now()
			// Обёртка над ResponseWriter, чтобы после обработки узнать статус ответа.
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK // обработчик ничего не записал явно
			}
			route := routePattern(r)

			httpRequestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(status)).Inc()
			metrics.ObserveWithTrace(r.Context(), httpRequestDuration.WithLabelValues(r.Method, route), time.Since(start).Seconds())
		})
	}
}

// routePattern возвращает шаблон маршрута (/test/videos-apify/{id}), а не реальный путь
// (/test/videos-apify/123): иначе каждый id породит отдельный временной ряд.
func routePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil || rctx.RoutePattern() == "" {
		return "unmatched"
	}
	return rctx.RoutePattern()
}
