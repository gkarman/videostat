package application

import (
	"context"

	"github.com/gkarman/demo/internal/domain/blogger"

)

type ApifyClient interface {
	Search(ctx context.Context, url string) (*blogger.Blogger, error)
}