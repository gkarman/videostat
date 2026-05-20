package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/notifier"
)

type VideoGenerationDoneHandler struct {
	notifier *notifier.TelegramNotifier
	log      *slog.Logger
}

func NewVideoGenerationDoneHandler(n *notifier.TelegramNotifier, log *slog.Logger) *VideoGenerationDoneHandler {
	return &VideoGenerationDoneHandler{notifier: n, log: log}
}

func (h *VideoGenerationDoneHandler) Handle(ctx context.Context, body []byte) error {
	var evt events.VideoGenerationDoneV1
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	text := fmt.Sprintf("✅ Видео готово\nID: %s\nS3: %s", evt.VideoID, evt.S3URL)
	if err := h.notifier.Notify(evt.ChatID, text); err != nil {
		h.log.Error("telegram notify failed", "err", err)
	}

	return nil
}
