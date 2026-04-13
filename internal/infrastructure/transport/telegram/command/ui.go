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
			tgbotapi.NewInlineKeyboardButtonData("Создать блогера", "create_blogger"),
			tgbotapi.NewInlineKeyboardButtonData("Список блогеров", "list_bloggers"),
			tgbotapi.NewInlineKeyboardButtonData("Список видео", "list_videos"),
			tgbotapi.NewInlineKeyboardButtonData("Экспорт видео в эксель", "export_videos"),
		),
	)
}

func (u *UI) CommandsText() string {
	return `<b>Доступные команды</b>

🚀 <code>/start</code>
📖 <code>/help</code>
➕ <code>/create_blogger</code>
➕ <code>/list_bloggers</code>
➕ <code>/export_videos</code>
📋 <code>/list_videos</code>`
}