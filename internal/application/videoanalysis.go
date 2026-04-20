package application

import "context"

type VideoAnalyzer interface {
	Analyze(ctx context.Context, videoURL string) ([]byte, error)
}
