package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/notifier"
)

type VideoCompositionDoneHandler struct {
	notifier *notifier.TelegramNotifier
	log      *slog.Logger
}

func NewVideoCompositionDoneHandler(n *notifier.TelegramNotifier, log *slog.Logger) *VideoCompositionDoneHandler {
	return &VideoCompositionDoneHandler{notifier: n, log: log}
}

func (h *VideoCompositionDoneHandler) Handle(ctx context.Context, body []byte) error {
	var evt events.VideoCompositionDoneV1
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	text := fmt.Sprintf("🎬 Итоговое видео готово\nИсходник: %s\nСмотреть: %s", evt.SourceURL, evt.ResultURL)
	if err := h.notifier.Notify(evt.ChatID, text); err != nil {
		h.log.Error("telegram notify failed", "err", err)
	}

	return nil
}
