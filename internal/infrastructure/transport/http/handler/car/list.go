package car

import (
	"net/http"

	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/response"
)

type ListHandler struct {
	service *service.List
}

func NewList(service *service.List) *ListHandler {
	return &ListHandler{
		service: service,
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	resp, err := h.service.Execute(r.Context(), requestdto.GetList{})
	if err != nil {
		log.Error("get cars failed", "error", err)
		response.ErrorJSON(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
