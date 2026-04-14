package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type BloggerCreatedHandler struct {
	command *command.FetchBloggerVideos
}

func NewBloggerCreatedHandler(log *slog.Logger, command *command.FetchBloggerVideos) *BloggerCreatedHandler {
	return &BloggerCreatedHandler{
		command: command,
	}
}

func (h *BloggerCreatedHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting BloggerCreated handler")

	var evt events.BloggerCreatedV1

	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	req := reqdto.FetchBloggerVideos{
		BloggerID: evt.BloggerID,
	}
	log.Debug("BloggerCreated handler", "req", req)

	err := h.command.Execute(ctx, req)
	if err != nil {
		log.Error("command error", "error", err)
	}

	return nil
}
