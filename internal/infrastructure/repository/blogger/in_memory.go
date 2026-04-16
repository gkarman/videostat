package blogger

import (
	"context"
	"sync"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type InMemoryRepo struct {
	mu              sync.Mutex
	bloggers        map[string]*blogger.Blogger // key = url
	videos          map[string]*blogger.Video
	SaveVideoErrFor map[string]error // externalID -> error
}

func NewInMemory() *InMemoryRepo {
	return &InMemoryRepo{
		bloggers: make(map[string]*blogger.Blogger),
		videos:   make(map[string]*blogger.Video),
		SaveVideoErrFor: make(map[string]error),
	}
}

func (r *InMemoryRepo) Save(_ context.Context, b *blogger.Blogger) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bloggers[b.URL] = b
	return nil
}

func (r *InMemoryRepo) ExistByUrl(_ context.Context, url string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.bloggers[url]
	return ok, nil
}

func (r *InMemoryRepo) GetById(ctx context.Context, id string) (*blogger.Blogger, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, b := range r.bloggers {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, blogger.ErrBloggerNotFound
}

func (r *InMemoryRepo) SaveVideo(ctx context.Context, video *blogger.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err, ok := r.SaveVideoErrFor[video.ExternalID]; ok {
		return err
	}

	r.videos[video.ExternalID] = video
	return nil
}

func (r *InMemoryRepo) ListVideosByBlogger(_ context.Context, bloggerID string) ([]*blogger.Video, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var res []*blogger.Video
	for _, v := range r.videos {
		if v.BloggerID == bloggerID {
			res = append(res, v)
		}
	}
	return res, nil
}
func (r *InMemoryRepo) List(ctx context.Context) ([]*blogger.Blogger, error) {
	return nil, nil
}
