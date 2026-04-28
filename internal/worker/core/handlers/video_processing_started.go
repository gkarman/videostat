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

type VideoProcessingStartedHandler struct {
	command *command.FetchVideoSources
}

func NewVideoProcessingStarted(log *slog.Logger, command *command.FetchVideoSources) *VideoProcessingStartedHandler {
	return &VideoProcessingStartedHandler{
		command: command,
	}
}

func (h *VideoProcessingStartedHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting VideoProcessingStartedHandler")

	var evt events.VideoProcessingStartedV1

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Debug("unmarshal", "body", string(body))
		return err
	}

	req := reqdto.FetchVideoSources{
		VideoID:  evt.VideoID,
		VideoURL: evt.VideoURL,
	}

	err := h.command.Execute(ctx, req)
	if err != nil {
		log.Error("command error", "error", err)
	}

	return nil
}
