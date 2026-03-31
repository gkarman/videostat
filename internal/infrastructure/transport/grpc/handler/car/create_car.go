package car

import (
	"context"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateCar(ctx context.Context, req *carv1.CreateCarRequest) (*carv1.CreateCarResponse, error) {
	reqDto := &requestdto.CreateCar{
		Name: req.GetName(),
	}

	respDto, err := h.createService.Execute(ctx, reqDto)
	if err != nil {
		h.log.Error("create car failed", "id", req.GetName(), "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &carv1.CreateCarResponse{
		Id: respDto.ID,
	}, nil
}
