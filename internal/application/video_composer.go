package application

import "context"

type CompositionClip struct {
	URL      string
	StartSec float64
	LenSec   float64
}

type CompositionRequest struct {
	AvatarURL string
	BrollClips []CompositionClip
	TotalSec   float64
}

type CompositionStatus struct {
	Status    string // processing | done | failed
	ResultURL string
	Error     string
}

type VideoComposer interface {
	Compose(ctx context.Context, req CompositionRequest) (externalID string, err error)
	GetStatus(ctx context.Context, externalID string) (*CompositionStatus, error)
}
