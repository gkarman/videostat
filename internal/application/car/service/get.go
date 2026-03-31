package service

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application/car/mapper"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/application/car/responsedto"
	"github.com/gkarman/demo/internal/domain/car"
)

type GetService struct {
	repo car.Repo
}

func NewGet(repo car.Repo) *GetService {
	return &GetService{
		repo: repo,
	}
}

func (s *GetService) Execute(ctx context.Context, req *requestdto.GetCar) (*responsedto.GetCar, error) {
	c, err := s.repo.GetByID(ctx, req.CarId)
	if err != nil {
		return nil, fmt.Errorf("GetService.Execute: %w", err)
	}

	return &responsedto.GetCar{
		Car: mapper.CarFromDomain(c),
	}, nil
}
