package application

import "context"

type VideoPromptGenerator interface {
	Generate(ctx context.Context, rawPayload []byte) (string, error)
	ProviderName() string
	ModelName() string
}
