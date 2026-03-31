package car

import (
	"context"
	"sync"

	"github.com/gkarman/demo/internal/domain/car"
)

type InMemoryRepo struct {
	mu   sync.RWMutex
	cars map[string]car.Car
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		cars: make(map[string]car.Car),
	}
}

func (r *InMemoryRepo) List(ctx context.Context) ([]*car.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*car.Car, 0, len(r.cars))
	for _, c := range r.cars {
		carCopy := c
		result = append(result, &carCopy)
	}

	return result, nil
}

func (r *InMemoryRepo) GetByID(ctx context.Context, id string) (*car.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.cars[id]
	if !ok {
		return nil, car.ErrNotFound
	}

	carCopy := c
	return &carCopy, nil
}

func (r *InMemoryRepo) Save(ctx context.Context, c *car.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[c.ID]; exists {
		return car.ErrAlreadyExists
	}

	r.cars[c.ID] = *c
	return nil
}

func (r *InMemoryRepo) Update(ctx context.Context, c *car.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[c.ID]; !exists {
		return car.ErrNotFound
	}

	r.cars[c.ID] = *c
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[id]; !exists {
		return car.ErrNotFound
	}

	delete(r.cars, id)
	return nil
}
