package application

import "context"

type Consumer interface {
	Consume(ctx context.Context, handler func([]byte) error) error
	Close() error
}
