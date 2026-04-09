package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/application/blogger/query"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token   string
	Debug   bool
	Timeout int
}

type userState struct {
	PlatformName string
	WaitingURL   bool
}

type Bot struct {
	botAPI       *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	log          *slog.Logger
	stop         chan struct{}
	done         chan struct{}

	createBlogger *command.CreateBlogger
	listBloggers *query.ListBloggers

	mu     sync.Mutex
	states map[int64]*userState
}

var defaultCommands = []tgbotapi.BotCommand{
	{Command: "start", Description: "Запуск бота"},
	{Command: "help", Description: "Список команд"},
	{Command: "create_blogger", Description: "Создать блогера"},
	{Command: "list_bloggers", Description: "Список всех блогеров"},
}

func NewBot(
	cfg *Config,
	log *slog.Logger,
	createBlogger *command.CreateBlogger,
	listBloggers *query.ListBloggers,
) (*Bot, error) {

	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("new telegram bot: %w", err)
	}

	botAPI.Debug = cfg.Debug

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = cfg.Timeout

	return &Bot{
		botAPI:        botAPI,
		updateConfig:  updateConfig,
		log:           log,
		stop:          make(chan struct{}),
		done:          make(chan struct{}),
		createBlogger: createBlogger,
		listBloggers: listBloggers,
		states:        make(map[int64]*userState),
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Info("start telegram bot")

	go func() {
		defer close(b.done)

		_ = b.setCommands()

		updates := b.botAPI.GetUpdatesChan(b.updateConfig)

		for {
			select {
			case <-ctx.Done():
				b.botAPI.StopReceivingUpdates()
				return

			case <-b.stop:
				b.botAPI.StopReceivingUpdates()
				return

			case update := <-updates:
				switch {
				case update.Message != nil:
					b.handleMessage(update.Message)

				case update.CallbackQuery != nil:
					b.handleCallback(update.CallbackQuery)
				}
			}
		}
	}()
}

func (b *Bot) Stop() {
	close(b.stop)
	<-b.done
}

func (b *Bot) setCommands() error {
	cfg := tgbotapi.NewSetMyCommands(defaultCommands...)
	_, err := b.botAPI.Request(cfg)
	return err
}

// -------------------- Message --------------------

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	// FSM: ждём URL
	if st, ok := b.getState(msg.Chat.ID); ok && st.WaitingURL {
		b.createBloggerFlow(msg, st)
		return
	}

	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}

	b.reply(msg.Chat.ID, "Я понимаю только команды. Нажми /help")
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {

	case "start":
		m := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Выбери действие:")
		m.ReplyMarkup = b.startInlineKeyboard()
		_, _ = b.botAPI.Send(m)

	case "help":
		b.reply(msg.Chat.ID, b.commandsText())

	case "create_blogger":
		b.askPlatform(msg.Chat.ID)

	case "list_bloggers":
		b.listBloggersFlow(msg.Chat.ID)

	default:
		b.reply(msg.Chat.ID, "Неизвестная команда")
	}
}

// -------------------- Callback --------------------

func (b *Bot) handleCallback(q *tgbotapi.CallbackQuery) {
	defer b.botAPI.Request(tgbotapi.NewCallback(q.ID, ""))

	switch {
	case q.Data == "create_blogger":
		b.askPlatform(q.Message.Chat.ID)

	case q.Data == "list_bloggers":
		b.listBloggersFlow(q.Message.Chat.ID)

	case strings.HasPrefix(q.Data, "platform_"):
		platform := strings.TrimPrefix(q.Data, "platform_")

		b.setState(q.Message.Chat.ID, &userState{
			PlatformName: platform,
			WaitingURL:   true,
		})

		b.reply(q.Message.Chat.ID, "Пришли ссылку на блогера")
	}
}

// -------------------- FSM Flow --------------------

func (b *Bot) createBloggerFlow(msg *tgbotapi.Message, st *userState) {
	ctx := context.Background()

	resp, err := b.createBlogger.Run(ctx, reqdto.CreateBlogger{
		URL:          msg.Text,
		PlatformName: st.PlatformName,
	})

	if err != nil {
		b.reply(msg.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	b.clearState(msg.Chat.ID)

	b.reply(msg.Chat.ID, fmt.Sprintf("Блогер создан ✅\nID: <code>%s</code>", resp.ID))
}

// -------------------- UI --------------------

func (b *Bot) askPlatform(chatID int64) {
	m := tgbotapi.NewMessage(chatID, "Выбери платформу:")
	m.ReplyMarkup = b.platformKeyboard()
	_, _ = b.botAPI.Send(m)
}

func (b *Bot) platformKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("youtube", "platform_youtube"),
			tgbotapi.NewInlineKeyboardButtonData("tiktok", "platform_tiktok"),
		),
	)
}

func (b *Bot) startInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать блогера", "create_blogger"),
			tgbotapi.NewInlineKeyboardButtonData("Список блогеров", "list_bloggers"),
		),
	)
}

func (b *Bot) commandsText() string {
	return `<b>Доступные команды</b>

🚀 <code>/start</code>
📖 <code>/help</code>
➕ <code>/create_blogger</code>`
}

func (b *Bot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, _ = b.botAPI.Send(msg)
}

// -------------------- State helpers --------------------

func (b *Bot) setState(chatID int64, st *userState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.states[chatID] = st
}

func (b *Bot) getState(chatID int64) (*userState, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	st, ok := b.states[chatID]
	return st, ok
}

func (b *Bot) clearState(chatID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.states, chatID)
}

func (b *Bot) listBloggersFlow(chatID int64) {
	ctx := context.Background()

	res, err := b.listBloggers.Run(ctx)
	if err != nil {
		b.reply(chatID, fmt.Sprintf("Ошибка при получении блогеров: %v", err))
		return
	}

	if len(res.Items) == 0 {
		b.reply(chatID, "Список блогеров пуст")
		return
	}

	var sb strings.Builder
	sb.WriteString("<b>Список блогеров:</b>\n\n")
	for i, bl := range res.Items {
		sb.WriteString(fmt.Sprintf("%d. <a href='%s'>%s</a> (%s)\n", i+1, bl.URL, bl.URL, bl.Platform))
	}

	b.reply(chatID, sb.String())
}
