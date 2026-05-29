package command

import (
	"context"
	"log/slog"
	"strings"

	appcmd "github.com/gkarman/demo/internal/application/blogger/command"
	appquery "github.com/gkarman/demo/internal/application/blogger/query"
	"github.com/gkarman/demo/internal/infrastructure/transport/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	log *slog.Logger

	sender           telegram.Sender
	allowedUsernames map[string]struct{}

	createBlogger     *appcmd.CreateBlogger
	listBloggersQuery *appquery.ListBloggers
	listVideosQuery   *appquery.ListVideos
	startProcessVideo *appcmd.StartProcessVideo
	resetVideoProcess *appcmd.ResetVideoProcess

	state *State
	ui    *UI
}

func NewRouter(
	log *slog.Logger,
	sender telegram.Sender,
	createBlogger *appcmd.CreateBlogger,
	listBloggers *appquery.ListBloggers,
	listVideos *appquery.ListVideos,
	startProcessVideo *appcmd.StartProcessVideo,
	resetVideoProcess *appcmd.ResetVideoProcess,
	allowedUsernames []string,
) *Router {
	allowed := make(map[string]struct{}, len(allowedUsernames))
	for _, u := range allowedUsernames {
		allowed[strings.TrimPrefix(strings.ToLower(u), "@")] = struct{}{}
	}
	return &Router{
		log:               log,
		sender:            sender,
		allowedUsernames:  allowed,
		createBlogger:     createBlogger,
		listBloggersQuery: listBloggers,
		listVideosQuery:   listVideos,
		startProcessVideo: startProcessVideo,
		resetVideoProcess: resetVideoProcess,
		state:             NewState(),
		ui:                NewUI(),
	}
}

func (r *Router) isAllowed(username string) bool {
	if len(r.allowedUsernames) == 0 {
		return true
	}
	_, ok := r.allowedUsernames[strings.ToLower(username)]
	return ok
}

func (r *Router) Commands() []tgbotapi.BotCommand {
	return commands()
}

func (r *Router) HandleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.From == nil || !r.isAllowed(msg.From.UserName) {
		r.send(msg.Chat.ID, "⛔ У вас нет доступа к этому боту.")
		return
	}

	if r.handleFSM(ctx, msg) {
		return
	}

	if msg.IsCommand() {
		r.handleCommand(ctx, msg)
	}
}

func (r *Router) HandleCallback(ctx context.Context, q *tgbotapi.CallbackQuery) {
	if q.From == nil || !r.isAllowed(q.From.UserName) {
		return
	}
	r.handleCallback(ctx, q)
}

func (r *Router) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_ = r.sender.Send(msg)
}

func (r *Router) sendWithKeyboard(chatID int64, text string, kb tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = kb
	_ = r.sender.Send(msg)
}
