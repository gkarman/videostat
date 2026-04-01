package video

import (
	"context"
)

type Repository interface {
	SaveSnapshot(ctx context.Context, snap *AccountSnapshot) error
	ExistsByPlatformAndExternalID(ctx context.Context, platformID int16, externalID string) (bool, error)
}
