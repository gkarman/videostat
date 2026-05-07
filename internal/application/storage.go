package application

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
}
