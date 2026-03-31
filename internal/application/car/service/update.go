package service

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application/car/mapper"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/application/car/responsedto"
	"github.com/gkarman/demo/internal/domain/car"
)

type UpdateService struct {
	repo car.Repo
}

func NewUpdate(repo car.Repo) *UpdateService {
	return &UpdateService{
		repo: repo,
	}
}

func (s *UpdateService) Execute(ctx context.Context, req *requestdto.UpdateCar) (*responsedto.UpdateCar, error) {
	if req.Name == "" {
		return nil, car.ErrEmptyName
	}

	c := car.New(req.CarId, req.Name)

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("UpdateService.Execute: %w", err)
	}

	return &responsedto.UpdateCar{
		Car: mapper.CarFromDomain(c),
	}, nil
}
