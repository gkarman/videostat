package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gkarman/demo/internal/infrastructure/logger"
)

func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithLogger(r.Context(), log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
