package videosourcesearcher

import (
	"context"
	"errors"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type InMemorySourceSearcher struct {
	Data map[string]string
	Errs map[string]error
}

func NewInMemorySourceSearcher() *InMemorySourceSearcher {
	return &InMemorySourceSearcher{
		Data: make(map[string]string),
		Errs: make(map[string]error),
	}
}

func (s *InMemorySourceSearcher) SearchUrl(ctx context.Context, v *blogger.Video) (string, error) {
	if err, ok := s.Errs[v.ID]; ok {
		return "", err
	}

	if url, ok := s.Data[v.ID]; ok {
		return url, nil
	}

	return "", errors.New("file source not found")
}