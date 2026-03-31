package car

import (
	"context"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetCar(ctx context.Context, req *carv1.GetCarRequest) (*carv1.GetCarResponse, error) {
	resp, err := h.getService.Execute(ctx, &requestdto.GetCar{CarId: req.GetId()})
	if err != nil {
		h.log.Error("get car failed", "id", req.GetId(), "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if resp.Car == nil {
		return nil, status.Error(codes.NotFound, "car not found")
	}

	return &carv1.GetCarResponse{
		Car: &carv1.CarInfo{
			Id:   resp.Car.ID,
			Name: resp.Car.Name,
		},
	}, nil
}
