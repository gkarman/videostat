package notifier

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramNotifier(token string) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}
	return &TelegramNotifier{bot: bot}, nil
}

func (n *TelegramNotifier) Notify(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := n.bot.Send(msg)
	return err
}
