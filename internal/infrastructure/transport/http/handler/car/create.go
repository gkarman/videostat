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
)

type CreateHandler struct {
	service *carservice.CreateService
}

func NewCreate(service *carservice.CreateService) *CreateHandler {
	return &CreateHandler{
		service: service,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req := &requestdto.CreateCar{
		Name: body.Name,
	}

	resp, err := h.service.Execute(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, car.ErrEmptyName):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		default:
			log.Error("save car failed", "error", err)
			response.ErrorJSON(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}
