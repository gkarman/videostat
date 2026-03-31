package application

import (
	"context"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, events []any)
}
