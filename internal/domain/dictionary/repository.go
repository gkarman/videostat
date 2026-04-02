package dictionary

import (
	"context"
)

type Repo interface {
	GetPlatformByName(context.Context, string) (*Platform, error)
}