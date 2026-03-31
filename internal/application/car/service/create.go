package service

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/application/car/responsedto"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/google/uuid"
)

type CreateService struct {
	repo       car.Repo
	dispatcher application.Dispatcher
}

func NewCreate(repo car.Repo, dispatcher application.Dispatcher) *CreateService {
	return &CreateService{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (s *CreateService) Execute(ctx context.Context, req *requestdto.CreateCar) (*responsedto.CreateCar, error) {
	if req.Name == "" {
		return nil, car.ErrEmptyName
	}

	id := uuid.New()
	c := car.New(id.String(), req.Name)

	if err := s.repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("CreateService.Execute: %w", err)
	}

	events := c.PullEvents()
	s.dispatcher.Dispatch(ctx, events)

	return &responsedto.CreateCar{
		ID: c.ID,
	}, nil
}
