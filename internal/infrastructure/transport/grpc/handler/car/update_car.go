package car

import (
	"context"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateCar(ctx context.Context, req *carv1.UpdateCarRequest) (*carv1.UpdateCarResponse, error) {
	resp, err := h.updateService.Execute(ctx, &requestdto.UpdateCar{
		CarId: req.GetId(),
		Name:  req.GetName(),
	})
	if err != nil {
		h.log.Error("update car failed", "id", req.GetId(), "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &carv1.UpdateCarResponse{
		Car: &carv1.CarInfo{
			Id:   resp.Car.ID,
			Name: resp.Car.Name,
		},
	}, nil
}
