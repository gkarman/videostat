package application

import "context"

type BrollGenerationStatus struct {
	Status   string // processing | done | failed
	VideoURL string
	Error    string
}

type BrollVideoGenerator interface {
	Submit(ctx context.Context, prompt string, durationSec int) (externalID string, err error)
	GetStatus(ctx context.Context, externalID string) (*BrollGenerationStatus, error)
}
