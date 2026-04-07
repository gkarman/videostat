package blogger

import (
	"context"
	"sync"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type InMemoryRepo struct {
	mu   sync.Mutex
	data map[string]*blogger.Blogger // key = url
}

func NewInMemory() *InMemoryRepo {
	return &InMemoryRepo{
		data: make(map[string]*blogger.Blogger),
	}
}

func (r *InMemoryRepo) Save(_ context.Context, b *blogger.Blogger) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[b.URL] = b
	return nil
}

func (r *InMemoryRepo) ExistByUrl(_ context.Context, url string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[url]
	return ok, nil
}