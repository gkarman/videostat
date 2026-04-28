package videoanalyzer

import (
	"context"
)

type FakeVideoAnalyzer struct {
}

func NewFakeVideoAnalyzer() *FakeVideoAnalyzer {
	return &FakeVideoAnalyzer{}
}

func (a *FakeVideoAnalyzer) ProviderName() string {
	return "fake"
}

func (a *FakeVideoAnalyzer) Analyze(ctx context.Context, videoURL string) ([]byte, error) {
	return nil, nil
}
