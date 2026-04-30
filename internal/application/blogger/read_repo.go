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

	Prompt      *string
	LLMProvider *string
	LLMModel    *string
}

type ReadRepo interface {
	ListBloggers(ctx context.Context) ([]BloggerRow, error)
	ListVideos(ctx context.Context) ([]VideoRow, error)
}
