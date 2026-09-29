package command

import (
	"context"
	"fmt"

	blogger_app "github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type FetchBloggerVideos struct {
	repo      blogger.Repo
	vSearcher blogger_app.VideoSearcher
}

func NewFetchBloggerVideos(repo blogger.Repo, vSearchers blogger_app.VideoSearcher) *FetchBloggerVideos {
	return &FetchBloggerVideos{
		repo:      repo,
		vSearcher: vSearchers,
	}
}

func (c *FetchBloggerVideos) Execute(ctx context.Context, req reqdto.FetchBloggerVideos) error {
	ctx = logger.WithField(ctx, logger.KeyBloggerID, req.BloggerID)
	log := logger.FromContext(ctx).With("component", "FetchBloggerVideos")

	b, err := c.repo.GetById(ctx, req.BloggerID)
	if err != nil {
		log.Error("get blogger by id failed", "error", err)
		return fmt.Errorf("get blogger by id: %w", err)
	}

	videos, err := c.vSearcher.Search(ctx, b)
	if err != nil {
		log.Error("video search failed", "error", err)
		return fmt.Errorf("video search: %w", err)
	}

	for _, v := range videos {
		log.Debug("video for save", "video", v)
		err := c.repo.SaveVideo(ctx, v)
		if err != nil {
			log.Error("save video failed", "external_video_id", v.ExternalID, "error", err)
			continue
		}
	}

	log.Info("videos fetched successfully", "count", len(videos))
	return nil
}
