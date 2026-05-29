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

func (u *UI) ConfirmRegenerateKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, перегенерировать", "regen_yes"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "regen_no"),
		),
	)
}

func (u *UI) StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1. Добавить ссылку на блогера", "create_blogger"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2. Экспорт видео в Excel", "export_videos"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3. Анализ видео по ссылке", "start_process_video"),
		),
	)
}
