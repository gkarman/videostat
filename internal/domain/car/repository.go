package car

import "context"

type Repo interface {
	List(ctx context.Context) ([]*Car, error)
	GetByID(ctx context.Context, id string) (*Car, error)
	Save(ctx context.Context, car *Car) error
	Update(ctx context.Context, car *Car) error
	Delete(ctx context.Context, id string) error
}