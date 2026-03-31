package car

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gkarman/demo/internal/application/car/requestdto"
	carservice "github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/response"
	"github.com/go-chi/chi/v5"
)

type UpdateHandler struct {
	service *carservice.UpdateService
}

func NewUpdate(service *carservice.UpdateService) *UpdateHandler {
	return &UpdateHandler{
		service: service,
	}
}

func (h *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id := chi.URLParam(r, "id")

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req := &requestdto.UpdateCar{
		CarId: id,
		Name:  body.Name,
	}

	resp, err := h.service.Execute(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, car.ErrEmptyName):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		default:
			log.Error("update car failed", "error", err)
			response.ErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
