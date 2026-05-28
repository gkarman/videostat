package application

import (
	"context"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type BrollGenerator interface {
	GenerateSegments(ctx context.Context, rawPayload []byte) ([]*blogger.BrollSegment, error)
}
