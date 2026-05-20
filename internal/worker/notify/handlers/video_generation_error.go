package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/notifier"
)

type VideoGenerationErrorHandler struct {
	notifier *notifier.TelegramNotifier
	log      *slog.Logger
}

func NewVideoGenerationErrorHandler(n *notifier.TelegramNotifier, log *slog.Logger) *VideoGenerationErrorHandler {
	return &VideoGenerationErrorHandler{notifier: n, log: log}
}

func (h *VideoGenerationErrorHandler) Handle(ctx context.Context, body []byte) error {
	var evt events.VideoGenerationErrorV1
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	text := fmt.Sprintf("❌ Ошибка генерации\nID: %s\nПричина: %s", evt.VideoID, evt.Reason)
	if err := h.notifier.Notify(evt.ChatID, text); err != nil {
		h.log.Error("telegram notify failed", "err", err)
	}

	return nil
}
