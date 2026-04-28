package application

import (
	"context"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type FileSourceSearcher interface {
	SearchUrl(ctx context.Context, v *blogger.Video) (string, error)
}
