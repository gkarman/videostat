package blogger

import (
	"context"
)

type Repo interface {
	Save(ctx context.Context, blogger *Blogger) error
	ExistByUrl(ctx context.Context, url string) (bool, error)
	GetById(ctx context.Context, id string) (*Blogger, error)
	SaveVideo(ctx context.Context, video *Video) error
	ListVideosByBlogger(ctx context.Context, bloggerID string) ([]*Video, error)
	List(ctx context.Context) ([]*Blogger, error)
	UpdateVideoStatus(ctx context.Context, videoID string, from VideoStatus, to VideoStatus) error
	GetVideoByUrl(ctx context.Context, url string) (*Video, error)
	GetVideoByID(ctx context.Context, id string) (*Video, error)
	UpdateVideoState(ctx context.Context, v *Video) error
	SaveVideoAnalysis(ctx context.Context, va *VideoAnalysis) error
	GetVideoAnalysisByVideoID(ctx context.Context, videoID string) (*VideoAnalysis, error)
	SaveVideoPrompt(ctx context.Context, vp *VideoPrompt) error
	GetVideoPromptByVideoID(ctx context.Context, videoID string) (*VideoPrompt, error)
	SaveVideoGeneration(ctx context.Context, vg *VideoGeneration) error
	GetVideoGenerationByExternalID(ctx context.Context, externalID string) (*VideoGeneration, error)
	ListPendingVideoGenerations(ctx context.Context) ([]*VideoGeneration, error)
	UpdateVideoGeneration(ctx context.Context, vg *VideoGeneration) error
}