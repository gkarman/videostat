package application

import (
	"context"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type VideoSearcher interface {
	Search(ctx context.Context, blogger *blogger.Blogger) ([]*blogger.Video, error)
}
