package platform

import (
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHTTPServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) *http.Server {
	router := http.NewRouter(log, db, d)
	return http.NewServer(
		log,
		router,
		http.Config{
			Addr:         cfg.ServerHttp.Addr,
			ReadTimeout:  time.Duration(cfg.ServerHttp.ReadTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(cfg.ServerHttp.WriteTimeoutSeconds) * time.Second,
		},
	)
}