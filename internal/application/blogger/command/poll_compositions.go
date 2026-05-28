package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/domain/blogger"
)

type PollCompositions struct {
	repo     blogger.Repo
	composer application.VideoComposer
}

func NewPollCompositions(repo blogger.Repo, composer application.VideoComposer) *PollCompositions {
	return &PollCompositions{repo: repo, composer: composer}
}

func (c *PollCompositions) Execute(ctx context.Context) error {
	compositions, err := c.repo.ListProcessingCompositions(ctx)
	if err != nil {
		return fmt.Errorf("list processing compositions: %w", err)
	}

	slog.Info("poll compositions", "count", len(compositions))

	for _, comp := range compositions {
		if comp.ExternalID == nil || *comp.ExternalID == "" {
			slog.Warn("composition has no external ID, skipping", "compositionID", comp.ID)
			continue
		}

		status, err := c.composer.GetStatus(ctx, *comp.ExternalID)
		if err != nil {
			slog.Error("get composition status", "compositionID", comp.ID, "externalID", *comp.ExternalID, "error", err)
			continue
		}
		slog.Info("composition status", "compositionID", comp.ID, "externalID", *comp.ExternalID, "status", status.Status)

		switch status.Status {
		case "done":
			comp.Status = blogger.CompositionStatusDone
			comp.ResultURL = &status.ResultURL
		case "failed":
			comp.Status = blogger.CompositionStatusFailed
			comp.Error = &status.Error
		default:
			continue
		}

		if err := c.repo.UpdateVideoComposition(ctx, comp); err != nil {
			return fmt.Errorf("update composition %s: %w", comp.ID, err)
		}
	}

	return nil
}
