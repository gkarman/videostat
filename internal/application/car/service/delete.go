package service

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/domain/car"
)

type DeleteService struct {
	repo car.Repo
}

func NewDelete(repo car.Repo) *DeleteService {
	return &DeleteService{
		repo: repo,
	}
}

func (s *DeleteService) Execute(ctx context.Context, req *requestdto.DeleteCar) error {
	if err := s.repo.Delete(ctx, req.CarId); err != nil {
		return fmt.Errorf("DeleteService.Execute: %w", err)
	}

	return nil
}
