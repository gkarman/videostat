package car

import (
	"errors"
	"net/http"

	"github.com/gkarman/demo/internal/application/car/requestdto"
	carservice "github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/response"
	"github.com/go-chi/chi/v5"
)

type GetByIDHandler struct {
	service *carservice.GetService
}

func NewGetCarHandler(service *carservice.GetService) *GetByIDHandler {
	return &GetByIDHandler{
		service: service,
	}
}

func (h *GetByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id := chi.URLParam(r, "id")

	resp, err := h.service.Execute(r.Context(), &requestdto.GetCar{CarId: id})
	if err != nil {
		switch {
		case errors.Is(err, car.ErrNotFound):
			response.ErrorJSON(w, http.StatusNotFound, err.Error())
		default:
			log.Error("get car failed", "error", err)
			response.ErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
