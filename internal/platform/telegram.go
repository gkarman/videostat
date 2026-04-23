package platform

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application/blogger/analytics"
	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/query"
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/repository/dictionary"
	"github.com/gkarman/demo/internal/infrastructure/transport/telegram"
	telergam_command "github.com/gkarman/demo/internal/infrastructure/transport/telegram/command"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTelegramBot(log *slog.Logger, cfg *config.Config, db *pgxpool.Pool, d *dispatcher.Dispatcher) (*telegram.Bot, error) {
	telegramCfg := &telegram.Config{
		Token:   cfg.TelegramBot.Token,
		Debug:   cfg.TelegramBot.Debug,
		Timeout: cfg.TelegramBot.Timeout,
	}

	repoBlogger := blogger.NewPostgres(db)
	repoDictionary := dictionary.NewPostgres(db)
	repoBloggerRead := blogger.NewQueryPostgres(db)



	createBloggerCmd := command.NewCreateBlogger(repoBlogger, repoDictionary, d)
	listBloggersQuery := query.NewListBloggers(repoBloggerRead)

	viralEnricher := analytics.NewViralEnricher()
	relevantEnricher := analytics.NewRelevanceEnricher()
	enricher := analytics.NewVideoMetricsEnricher(viralEnricher, relevantEnricher)
	listVideosQuery := query.NewListVideos(repoBloggerRead, enricher)

	startProcessVideoCmd := command.NewStartProcessVideo(repoBlogger, d)

	bot, err := telegram.NewBot(telegramCfg, log)
	if err != nil {
		return nil, err
	}

	router := telergam_command.NewRouter(
		log,
		bot,
		createBloggerCmd,
		listBloggersQuery,
		listVideosQuery,
		startProcessVideoCmd,
	)

	bot.SetHandler(router)

	return bot, nil
}
