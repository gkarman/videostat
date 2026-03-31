package service

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application/car/mapper"
	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/application/car/responsedto"
	"github.com/gkarman/demo/internal/domain/car"
)

type List struct {
	repo car.Repo
}

func NewList(repo car.Repo) *List {
	return &List{
		repo: repo,
	}
}

func (s *List) Execute(ctx context.Context, _ requestdto.GetList) (*responsedto.GetList, error) {
	cars, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf(`List.handel: %w`, err)
	}
	resp := &responsedto.GetList{
		Cars: mapper.CarsFromDomain(cars),
	}
	return resp, nil
}
