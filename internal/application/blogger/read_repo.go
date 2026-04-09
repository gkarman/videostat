package blogger

import (
	"context"
)

// BloggerRow — сырые данные из репозитория (infrastructure)
type BloggerRow struct {
	ID       string
	URL      string
	Platform string
}

type ReadRepo interface {
	List(ctx context.Context) ([]BloggerRow, error)
}