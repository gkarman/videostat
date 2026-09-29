package platform

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/infrastructure/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer поднимает HTTP-сервер, который отдаёт только /metrics,
// и останавливает его, когда завершается ctx. Нужен воркерам: своего HTTP у них нет.
func StartMetricsServer(ctx context.Context, log *slog.Logger, addr string) {
	server := http.NewServer(
		log.With("component", "metrics"),
		promhttp.Handler(),
		http.Config{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)
	server.Start()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Stop(shutdownCtx); err != nil {
			log.Error("stop metrics server", "error", err)
		}
	}()
}
