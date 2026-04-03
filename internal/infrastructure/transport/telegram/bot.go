package telegram

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token   string
	Debug   bool
	Timeout int
}

type Bot struct {
	botAPI       *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	log          *slog.Logger
	stop         chan struct{}
	done         chan struct{}
}

func NewBot(cfg *Config, log *slog.Logger) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("new telegram bot: %w", err)
	}

	botAPI.Debug = cfg.Debug
	log.Info("telegram authorized on account %s", botAPI.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = cfg.Timeout

	return &Bot{
		botAPI:       botAPI,
		updateConfig: updateConfig,
		log:          log,
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Info("start telegram bot")
	go func() {
		defer close(b.done)

		updates := b.botAPI.GetUpdatesChan(b.updateConfig)

		for {
			select {
			case <-ctx.Done():
				b.log.Info("telegram ctx done")
				b.botAPI.StopReceivingUpdates()
				return

			case <-b.stop:
				b.log.Info("telegram stop signal")
				b.botAPI.StopReceivingUpdates()
				return

			case update := <-updates:
				if update.Message == nil {
					continue
				}
				b.handle(update)
			}
		}
	}()
}

func (b *Bot) handle(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Ты написал: "+update.Message.Text)

	if _, err := b.botAPI.Send(msg); err != nil {
		b.log.Error("send message", "err", err)
	}
}

func (b *Bot) Stop() {
	b.log.Info("stop telegram bot")
	close(b.stop)
	<-b.done
}
