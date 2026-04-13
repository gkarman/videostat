package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "help", Description: "Список команд"},
		{Command: "create_blogger", Description: "Создать блогера"},
		{Command: "list_bloggers", Description: "Список блогеров"},
	}
}

func (r *Router) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		r.sendWithKeyboard(
			msg.Chat.ID,
			"Привет! Выбери действие:",
			r.ui.StartKeyboard(),
		)
	case "help":
		r.send(msg.Chat.ID, r.ui.CommandsText())
	case "create_blogger":
		r.askPlatform(msg.Chat.ID)
	case "list_bloggers":
		r.listBloggers(ctx, msg.Chat.ID)
	default:
		r.send(msg.Chat.ID, "Неизвестная команда")
	}
}
