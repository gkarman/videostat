package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) PlatformKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("youtube", "platform_youtube"),
			tgbotapi.NewInlineKeyboardButtonData("tiktok", "platform_tiktok"),
			tgbotapi.NewInlineKeyboardButtonData("instagram", "platform_instagram"),
		),
	)
}

func (u *UI) StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1. Добавить ссылку на блогера", "create_blogger"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2. Список блогеров", "list_bloggers"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3. Экспорт видео в Excel", "export_videos"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4. Анализ видео по ссылке", "start_process_video"),
		),
	)
}
