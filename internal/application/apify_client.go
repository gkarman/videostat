package application

import (
	"context"

	"github.com/gkarman/demo/internal/domain/video"
)

type ApifyClient interface {
	FetchAccountSnapshot(ctx context.Context, url string) (*video.AccountSnapshot, error)
}