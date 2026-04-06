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

var defaultCommands = []tgbotapi.BotCommand{
	{Command: "start", Description: "Запуск бота"},
	{Command: "help", Description: "Список команд"},
}

func NewBot(cfg *Config, log *slog.Logger) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("new telegram bot: %w", err)
	}

	botAPI.Debug = cfg.Debug

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
		if err := b.setCommands(); err != nil {
			b.log.Error("set commands", "err", err)
		}
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
	msg := update.Message

	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}

	b.reply(msg.Chat.ID, "Я понимаю только команды. Нажми /help")
}

func (b *Bot) Stop() {
	b.log.Info("stop telegram bot")
	close(b.stop)
	<-b.done
}

func (b *Bot) setCommands() error {
	cfg := tgbotapi.NewSetMyCommands(defaultCommands...)
	_, err := b.botAPI.Request(cfg)
	return err
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {

	case "start":
		//msg := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Выбери команду из меню ниже:")
		//msg.ReplyMarkup = b.startKeyboard()
		//msg.ParseMode = tgbotapi.ModeHTML
		//if _, err := b.botAPI.Send(msg); err != nil {
		//	b.log.Error("send message", "err", err)
		//}
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Выбери команду:")
		msg.ReplyMarkup = b.startInlineKeyboard()
		b.botAPI.Send(msg)

	case "help":
		b.reply(msg.Chat.ID, b.commandsText())


	default:
		b.reply(msg.Chat.ID, "Неизвестная команда")
	}
}

func (b *Bot) commandsText() string {
	return `<b>Доступные команды</b>

🚀 <code>/start</code> — запуск бота
📖 <code>/help</code> — список команд

Выбери команду из меню или введи её вручную.`
}

func (b *Bot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := b.botAPI.Send(msg); err != nil {
		b.log.Error("send message", "err", err)
	}
}


func (b *Bot) startKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
			tgbotapi.NewKeyboardButton("/help"),
		),
	)
}

func (b *Bot) startInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Start", "cmd_start"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "cmd_help"),
		),
	)
}