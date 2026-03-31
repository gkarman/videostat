package platform

import (
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)
	return log
}
