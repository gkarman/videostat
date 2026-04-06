package platform

import (
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/repository/dictionary"
	"github.com/gkarman/demo/internal/infrastructure/transport/telegram"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTelegramBot(log *slog.Logger, cfg *config.Config, db *pgxpool.Pool, d *dispatcher.Dispatcher) (*telegram.Bot, error) {
	telegramCfg := &telegram.Config{
		Token: cfg.TelegramBot.Token,
		Debug: cfg.TelegramBot.Debug,
		Timeout: cfg.TelegramBot.Timeout,
	}

	repoBlogger := blogger.NewPostgres(db)
	repoDictionary := dictionary.NewPostgres(db)
	createBlogerCmd := command.NewCreateBlogger(repoBlogger, repoDictionary, d)

	bot, err := telegram.NewBot(telegramCfg, log, createBlogerCmd)
	if err != nil {
		return  nil, fmt.Errorf("error creating telegram bot: %w", err)
	}

	return bot, nil
}
