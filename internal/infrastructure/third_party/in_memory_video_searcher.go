package third_party

import (
	"context"
	"time"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type InMemoryVideoSearcher struct {
	Videos []*blogger.Video
	Err    error
}

func NewInMemoryVideoSearcher(videos []*blogger.Video, err error) *InMemoryVideoSearcher {
	return &InMemoryVideoSearcher{
		Videos: videos,
		Err:    err,
	}
}

func (f *InMemoryVideoSearcher) Search(_ context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	now := time.Now()
	for _, v := range f.Videos {
		v.BloggerID = b.ID
		v.CreatedAt = now
	}

	return f.Videos, nil
}
