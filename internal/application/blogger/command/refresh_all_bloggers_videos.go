package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type RefreshAllBloggers struct {
	repo  blogger.Repo
	fetch *FetchBloggerVideos
}

func NewRefreshAllBloggers(repo blogger.Repo, fetch *FetchBloggerVideos) *RefreshAllBloggers {
	return &RefreshAllBloggers{
		repo:  repo,
		fetch: fetch,
	}
}

func (c *RefreshAllBloggers) Execute(ctx context.Context) error {
	log := logger.FromContext(ctx).With("component", "RefreshAllBloggers")

	bloggers, err := c.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("list bloggers: %w", err)
	}

	for _, b := range bloggers {
		err := c.fetch.Execute(ctx, reqdto.FetchBloggerVideos{
			BloggerID: b.ID,
		})
		if err != nil {
			log.Error("fetch failed", "bloggerID", b.ID, "err", err)
			continue
		}
	}

	log.Info("refresh finished", "count", len(bloggers))
	return nil
}
