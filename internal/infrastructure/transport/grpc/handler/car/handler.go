package car

import (
	"log/slog"

	grpcCarv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/service"
)

type Handler struct {
	grpcCarv1.UnimplementedCarServer
	log           *slog.Logger
	getService    *service.GetService
	listService   *service.List
	createService *service.CreateService
	updateService *service.UpdateService
	deleteService *service.DeleteService
}

func New(
	log *slog.Logger,
	getService *service.GetService,
	listService *service.List,
	createService *service.CreateService,
	updateService *service.UpdateService,
	deleteService *service.DeleteService,
) *Handler {
	return &Handler{
		log:           log,
		getService:    getService,
		listService:   listService,
		createService: createService,
		updateService: updateService,
		deleteService: deleteService,
	}
}
