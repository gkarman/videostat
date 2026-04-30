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

type VideoAnalyzeDoneHandler struct {
	command *command.GenerateVideoPrompt
}

func NewVideoAnalyzeDoneHandler(log *slog.Logger, command *command.GenerateVideoPrompt) *VideoAnalyzeDoneHandler {
	return &VideoAnalyzeDoneHandler{command: command}
}

func (h *VideoAnalyzeDoneHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting VideoAnalyzeDoneHandler")

	var evt events.VideoAnalyzeDoneV1

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Debug("unmarshal failed", "body", string(body))
		return err
	}

	err := h.command.Run(ctx, reqdto.GenerateVideoPrompt{VideoID: evt.VideoID})
	if err != nil {
		log.Error("command error", "error", err)
	}

	return nil
}
