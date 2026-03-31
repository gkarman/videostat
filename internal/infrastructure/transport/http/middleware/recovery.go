package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gkarman/demo/internal/infrastructure/logger"
)

const maxStackLogBytes = 8 * 1024

func Recovery() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log := logger.FromContext(r.Context())
					stack, truncated := limitedStack()
					log.Error(
						"http panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"panic", fmt.Sprint(rec),
						"stack", stack,
						"stack_truncated", truncated,
					)
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func limitedStack() (string, bool) {
	stack := debug.Stack()
	if len(stack) <= maxStackLogBytes {
		return string(stack), false
	}
	return string(stack[:maxStackLogBytes]), true
}
