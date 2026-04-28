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

type VideoSourceFoundHandler struct {
	command *command.AnalyzeVideo
}

func NewVideoSourceFoundHandler(log *slog.Logger, command *command.AnalyzeVideo) *VideoSourceFoundHandler {
	return &VideoSourceFoundHandler{command: command}
}

func (h *VideoSourceFoundHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting VideoSourceFoundHandler")

	var evt events.VideoSourceFoundV1

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Debug("unmarshal failed", "body", string(body))
		return err
	}

	req := reqdto.AnalyzeVideo{
		VideoID: evt.VideoID,
		FileURL: evt.FileURL,
	}

	err := h.command.Run(ctx, req)
	if err != nil {
		log.Error("command error", "error", err)
	}

	return nil
}
