package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type SubmitVideoGeneration struct {
	r blogger.Repo
	g application.VideoGenerator
}

func NewSubmitVideoGeneration(r blogger.Repo, g application.VideoGenerator) *SubmitVideoGeneration {
	return &SubmitVideoGeneration{r: r, g: g}
}

func (c *SubmitVideoGeneration) Run(ctx context.Context, req reqdto.SubmitVideoGeneration) error {
	log := logger.FromContext(ctx).With("component", "SubmitVideoGeneration", "videoID", req.VideoID)

	vp, err := c.r.GetVideoPromptByVideoID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video prompt: %w", err)
	}

	externalID, err := c.g.Submit(ctx, vp.Script)
	if err != nil {
		return fmt.Errorf("submit to %s: %w", c.g.Platform(), err)
	}

	vg := &blogger.VideoGeneration{
		ID:         uuid.NewString(),
		VideoID:    req.VideoID,
		Platform:   c.g.Platform(),
		ExternalID: externalID,
		Status:     blogger.VideoGenerationPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err = c.r.SaveVideoGeneration(ctx, vg); err != nil {
		return fmt.Errorf("save video generation: %w", err)
	}

	log.Info("video generation submitted", "platform", vg.Platform, "externalID", externalID)
	return nil
}
