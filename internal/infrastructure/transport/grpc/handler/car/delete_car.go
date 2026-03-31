package car

import (
	"context"
	"errors"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/domain/car"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) DeleteCar(ctx context.Context, req *carv1.DeleteCarRequest) (*carv1.DeleteCarResponse, error) {
	err := h.deleteService.Execute(ctx, &requestdto.DeleteCar{CarId: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, car.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			h.log.Error("delete car failed", "id", req.GetId(), "error", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &carv1.DeleteCarResponse{}, nil
}
