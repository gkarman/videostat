package application

import "context"

type GenerationStatus struct {
	Status       string
	VideoURL     string
	ErrorMessage string
}

type VideoGenerator interface {
	Platform() string
	Submit(ctx context.Context, script string) (externalID string, err error)
	GetStatus(ctx context.Context, externalID string) (*GenerationStatus, error)
}
