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
	log            *slog.Logger
	submitGen      *command.SubmitVideoGeneration
	generateBroll  *command.GenerateBrollSegments
}

func NewVideoAnalyzeDoneHandler(log *slog.Logger, submitGen *command.SubmitVideoGeneration, generateBroll *command.GenerateBrollSegments) *VideoAnalyzeDoneHandler {
	return &VideoAnalyzeDoneHandler{log: log, submitGen: submitGen, generateBroll: generateBroll}
}

func (h *VideoAnalyzeDoneHandler) Handle(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Debug("starting VideoAnalyzeDoneHandler")

	var evt events.VideoAnalyzeDoneV1

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Debug("unmarshal failed", "body", string(body))
		return err
	}

	if err := h.submitGen.Run(ctx, reqdto.SubmitVideoGeneration{VideoID: evt.VideoID}); err != nil {
		log.Error("submit video generation failed", "error", err)
	}

	if err := h.generateBroll.Run(ctx, reqdto.GenerateBrollSegments{VideoID: evt.VideoID}); err != nil {
		h.log.Error("generate broll segments failed", "error", err)
	}

	return nil
}
