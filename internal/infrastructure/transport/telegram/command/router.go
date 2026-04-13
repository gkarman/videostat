package command

import (
	"context"
	"log/slog"

	appcmd "github.com/gkarman/demo/internal/application/blogger/command"
	appquery "github.com/gkarman/demo/internal/application/blogger/query"
	"github.com/gkarman/demo/internal/infrastructure/transport/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	log *slog.Logger

	sender telegram.Sender

	createBlogger *appcmd.CreateBlogger
	listBloggersQuery  *appquery.ListBloggers

	state *State
	ui    *UI
}

func NewRouter(
	log *slog.Logger,
	sender telegram.Sender,
	createBlogger *appcmd.CreateBlogger,
	listBloggers *appquery.ListBloggers,
) *Router {
	return &Router{
		log:           log,
		sender:        sender,
		createBlogger: createBlogger,
		listBloggersQuery:  listBloggers,
		state:         NewState(),
		ui:            NewUI(),
	}
}

func (r *Router) Commands() []tgbotapi.BotCommand {
	return commands()
}

func (r *Router) HandleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if r.handleFSM(ctx, msg) {
		return
	}

	if msg.IsCommand() {
		r.handleCommand(ctx, msg)
	}
}

func (r *Router) HandleCallback(ctx context.Context, q *tgbotapi.CallbackQuery) {
	r.handleCallback(ctx, q)
}

func (r *Router) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_ = r.sender.Send(msg)
}
