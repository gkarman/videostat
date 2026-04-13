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

type Handler interface {
	HandleMessage(ctx context.Context, msg *tgbotapi.Message)
	HandleCallback(ctx context.Context, q *tgbotapi.CallbackQuery)
	Commands() []tgbotapi.BotCommand
}

type Sender interface {
	Send(msg tgbotapi.Chattable) error
}

type Bot struct {
	api     *tgbotapi.BotAPI
	cfg     *Config
	log     *slog.Logger
	handler Handler

	stop chan struct{}
	done chan struct{}
}

func NewBot(cfg *Config, log *slog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("new telegram bot: %w", err)
	}

	api.Debug = cfg.Debug

	return &Bot{
		api:     api,
		cfg:     cfg,
		log:     log,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Info("start telegram bot")

	go func() {
		defer close(b.done)

		b.setCommands()

		updateCfg := tgbotapi.NewUpdate(0)
		updateCfg.Timeout = b.cfg.Timeout

		updates := b.api.GetUpdatesChan(updateCfg)

		for {
			select {
			case <-ctx.Done():
				b.api.StopReceivingUpdates()
				return

			case <-b.stop:
				b.api.StopReceivingUpdates()
				return

			case upd := <-updates:
				switch {
				case upd.Message != nil:
					b.handler.HandleMessage(ctx, upd.Message)

				case upd.CallbackQuery != nil:
					b.handler.HandleCallback(ctx, upd.CallbackQuery)
				}
			}
		}
	}()
}

func (b *Bot) Stop() {
	close(b.stop)
	<-b.done
}

func (b *Bot) setCommands() {
	cmds := b.handler.Commands()
	if len(cmds) == 0 {
		return
	}

	cfg := tgbotapi.NewSetMyCommands(cmds...)
	_, err := b.api.Request(cfg)
	if err != nil {
		b.log.Error("set telegram commands", slog.Any("err", err))
	}
}

func (b *Bot) Send(msg tgbotapi.Chattable) error {
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) SetHandler(h Handler) {
	b.handler = h
}