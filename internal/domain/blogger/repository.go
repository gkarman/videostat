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
}