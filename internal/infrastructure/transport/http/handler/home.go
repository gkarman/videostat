package handler

import (
	"net/http"

	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hello World")); err != nil {
		log.Error("write response", "err", err)
	}
}
