package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
)

type CarUpdatedHandler struct {
	log *slog.Logger
}

func NewCarUpdatedHandler(log *slog.Logger) *CarUpdatedHandler {
	return &CarUpdatedHandler{log: log}
}

func (h *CarUpdatedHandler) Handle(body []byte) error {
	var evt events.CarUpdatedV1

	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}

	h.log.Info("car updated",
		"id", evt.CarID,
		"name", evt.NameOld,
	)

	return nil
}