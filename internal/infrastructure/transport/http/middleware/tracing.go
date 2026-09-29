package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Tracing создаёт спан на каждый входящий HTTP-запрос. Если во входящем запросе есть
// заголовок traceparent (запрос пришёл из другого сервиса) — спан продолжит тот же трейс.
// Спан кладётся в r.Context(): всё, что вызывается дальше с этим ctx (SQL, внешние API),
// станет его дочерними спанами.
func Tracing() func(next http.Handler) http.Handler {
	return otelhttp.NewMiddleware("http.server",
		// /metrics дёргает Prometheus каждые 15 секунд — такие трейсы только мусорят.
		otelhttp.WithFilter(func(r *http.Request) bool { return r.URL.Path != "/metrics" }),
	)
}

// NameSpanByRoute называет спан по шаблону маршрута: "GET /test/videos-apify/{id}".
// Шаблон известен только после роутинга, поэтому имя задаётся после next.ServeHTTP.
// Как и в метриках, берём шаблон, а не путь с id — иначе поиск по имени спана бесполезен.
func NameSpanByRoute() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			trace.SpanFromContext(r.Context()).SetName(r.Method + " " + routePattern(r))
		})
	}
}
