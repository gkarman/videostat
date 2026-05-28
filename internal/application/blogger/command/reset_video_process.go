package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type ResetVideoProcess struct {
	repo blogger.Repo
}

func NewResetVideoProcess(repo blogger.Repo) *ResetVideoProcess {
	return &ResetVideoProcess{repo: repo}
}

func (c *ResetVideoProcess) Run(ctx context.Context, videoID string) error {
	if err := c.repo.DeleteBrollSegmentsByVideoID(ctx, videoID); err != nil {
		return fmt.Errorf("delete broll segments: %w", err)
	}
	if err := c.repo.DeleteCompositionsByVideoID(ctx, videoID); err != nil {
		return fmt.Errorf("delete compositions: %w", err)
	}
	if err := c.repo.DeleteVideoGenerationsByVideoID(ctx, videoID); err != nil {
		return fmt.Errorf("delete video generations: %w", err)
	}
	if err := c.repo.DeleteVideoWatchersByVideoID(ctx, videoID); err != nil {
		return fmt.Errorf("delete video watchers: %w", err)
	}
	return nil
}
