package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
)

type PollBrollGenerations struct {
	repo    blogger.Repo
	gen     application.BrollVideoGenerator
	compose *ComposeFinalVideo
}

func NewPollBrollGenerations(repo blogger.Repo, gen application.BrollVideoGenerator, compose *ComposeFinalVideo) *PollBrollGenerations {
	return &PollBrollGenerations{repo: repo, gen: gen, compose: compose}
}

func (c *PollBrollGenerations) Execute(ctx context.Context) error {
	segments, err := c.repo.ListProcessingBrollSegments(ctx)
	if err != nil {
		return fmt.Errorf("list processing broll segments: %w", err)
	}

	slog.Info("poll broll generations", "count", len(segments))

	finishedVideoIDs := map[string]struct{}{}

	for _, s := range segments {
		if s.GenerationExternalID == nil {
			continue
		}

		status, err := c.gen.GetStatus(ctx, *s.GenerationExternalID)
		if err != nil {
			slog.Error("get broll status", "segmentID", s.ID, "externalID", *s.GenerationExternalID, "error", err)
			continue
		}

		slog.Info("broll segment status", "segmentID", s.ID, "externalID", *s.GenerationExternalID, "status", status.Status)

		switch status.Status {
		case "done":
			s.GenerationStatus = blogger.BrollStatusDone
			s.BrollURL = &status.VideoURL
		case "failed":
			s.GenerationStatus = blogger.BrollStatusFailed
			s.GenerationError = &status.Error
		default:
			continue
		}

		if err := c.repo.UpdateBrollSegment(ctx, s); err != nil {
			return fmt.Errorf("update broll segment %s: %w", s.ID, err)
		}

		finishedVideoIDs[s.VideoID] = struct{}{}
	}

	for videoID := range finishedVideoIDs {
		active, err := c.repo.CountActiveBrollSegments(ctx, videoID)
		if err != nil || active > 0 {
			continue
		}
		slog.Info("all broll segments done, composing final video", "videoID", videoID)
		if err := c.compose.Run(ctx, reqdto.ComposeFinalVideo{VideoID: videoID}); err != nil {
			if errors.Is(err, ErrAvatarNotReady) {
				slog.Info("avatar not ready yet, will retry when heygen finishes", "videoID", videoID)
				continue
			}
			return fmt.Errorf("compose final video %s: %w", videoID, err)
		}
	}

	return nil
}
