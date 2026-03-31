package car

import (
	"context"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetCarList(ctx context.Context, _ *carv1.GetCarListRequest) (*carv1.GetCarListResponse, error) {
	resp, err := h.listService.Execute(ctx, requestdto.GetList{})
	if err != nil {
		h.log.Error("get car list failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	cars := make([]*carv1.CarInfo, 0, len(resp.Cars))
	for _, c := range resp.Cars {
		cars = append(cars, &carv1.CarInfo{
			Id:   c.ID,
			Name: c.Name,
		})
	}

	return &carv1.GetCarListResponse{
		Cars: cars,
	}, nil
}
