package blogger

import (
	"context"
	"time"
)

type BloggerRow struct {
	ID       string
	URL      string
	Platform string
}

type VideoRow struct {
	ID          string
	Platform    string
	BloggerURL  string
	URL         string
	Title       string
	Views       int
	Likes       int
	Comments    int
	PublishedAt time.Time
	CreatedAt   time.Time

	Status       string
	ErrorStage   *string
	ErrorMessage *string

	AnalysisProvider *string
	RawPayload       *string

	Brief       *string
	Script      *string
	LLMProvider *string
	LLMModel    *string

	GenerationPlatform      *string
	GenerationExternalID    *string
	GenerationStatus        *string
	GenerationS3URL         *string
	GenerationErrorMessage  *string
}

type ReadRepo interface {
	ListBloggers(ctx context.Context) ([]BloggerRow, error)
	ListVideos(ctx context.Context) ([]VideoRow, error)
}
