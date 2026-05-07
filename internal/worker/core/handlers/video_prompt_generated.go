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

type VideoPromptGeneratedHandler struct {
	command *command.SubmitVideoGeneration
}

func NewVideoPromptGeneratedHandler(log *slog.Logger, command *command.SubmitVideoGeneration) *VideoPromptGeneratedHandler {
	return &VideoPromptGeneratedHandler{command: command}
}

func (h *VideoPromptGeneratedHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting VideoPromptGeneratedHandler")

	var evt events.VideoPromptGeneratedV1
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	if err := h.command.Run(ctx, reqdto.SubmitVideoGeneration{VideoID: evt.VideoID}); err != nil {
		log.Error("command error", "error", err)
	}

	return nil
}
