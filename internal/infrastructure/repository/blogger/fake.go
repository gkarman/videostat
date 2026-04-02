package blogger

import (
	"context"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type FakeRepo struct{}

func NewFake() *FakeRepo {
	return &FakeRepo{}
}

func (r *FakeRepo) Save(_ context.Context, _ *blogger.Blogger) error {
	return nil
}
