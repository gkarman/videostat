package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "start", Description: "Запуск бота"},
		{Command: "create_blogger", Description: "Создать блогера"},
		{Command: "list_bloggers", Description: "Список блогеров"},
		{Command: "list_videos", Description: "Список видео"},
		{Command: "export_videos", Description: "Экспорт видео"},
		{Command: "start_process_video", Description: "Обработать видео по ссылке"},
	}
}

func (r *Router) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	r.log.Debug("command from telegram", "command", msg.Command())

	switch msg.Command() {
	case "start":
		r.sendWithKeyboard(
			msg.Chat.ID,
			"Привет! Выбери действие:",
			r.ui.StartKeyboard(),
		)
	case "create_blogger":
		r.askPlatform(msg.Chat.ID)
	case "list_bloggers":
		r.listBloggers(ctx, msg.Chat.ID)
	case "list_videos":
		r.listVideos(ctx, msg.Chat.ID)
	case "export_videos":
		r.exportVideos(ctx, msg.Chat.ID)
	case "start_process_video":
		r.state.Set(msg.Chat.ID, &userState{
			WaitingVideoURL: true,
		})
		r.send(msg.Chat.ID, "Пришлите ссылку на видео")
	default:
		r.send(msg.Chat.ID, "Неизвестная команда")
	}
}
