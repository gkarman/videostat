package platform

import (
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

// NewLogger создаёт логгер сервиса. service — имя файла логов и метка service в Loki.
func NewLogger(cfg *config.Config, service string) *slog.Logger {
	log := logger.New(logger.Config{
		Level:   cfg.Logger.Level,
		Dir:     cfg.Logger.Dir,
		Service: service,
	})
	slog.SetDefault(log)
	return log
}
