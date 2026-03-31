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
	_ = logger.FromContext(r.Context()) // если понадобится логировать
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
