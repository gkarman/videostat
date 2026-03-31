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

type DeleteHandler struct {
	service *carservice.DeleteService
}

func NewDelete(service *carservice.DeleteService) *DeleteHandler {
	return &DeleteHandler{
		service: service,
	}
}

func (h *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id := chi.URLParam(r, "id")

	err := h.service.Execute(r.Context(), &requestdto.DeleteCar{CarId: id})
	if err != nil {
		switch {
		case errors.Is(err, car.ErrNotFound):
			response.ErrorJSON(w, http.StatusNotFound, err.Error())
		default:
			log.Error("delete car failed", "error", err)
			response.ErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
