package platform

import (
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/transport/telegram"
)

func NewTelegramBot(log *slog.Logger, cfg *config.Config) (*telegram.Bot, error) {
	telegramCfg := &telegram.Config{
		Token: cfg.TelegramBot.Token,
		Debug: cfg.TelegramBot.Debug,
		Timeout: cfg.TelegramBot.Timeout,
	}

	bot, err := telegram.NewBot(telegramCfg, log)
	if err != nil {
		return  nil, fmt.Errorf("error creating telegram bot: %w", err)
	}

	return bot, nil
}
