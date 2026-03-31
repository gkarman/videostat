package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
)

type CarCreatedHandler struct {
	log *slog.Logger
}

func NewCarCreatedHandler(log *slog.Logger) *CarCreatedHandler {
	return &CarCreatedHandler{log: log}
}

func (h *CarCreatedHandler) Handle(body []byte) error {
	var evt events.CarCreatedV1

	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	h.log.Info("car created",
		"id", evt.CarID,
		"name", evt.Name,
	)

	return nil
}