package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type PollCompositions struct {
	repo       blogger.Repo
	composer   application.VideoComposer
	dispatcher application.Dispatcher
}

func NewPollCompositions(repo blogger.Repo, composer application.VideoComposer, dispatcher application.Dispatcher) *PollCompositions {
	return &PollCompositions{repo: repo, composer: composer, dispatcher: dispatcher}
}

func (c *PollCompositions) Execute(ctx context.Context) error {
	log := logger.FromContext(ctx).With("component", "PollCompositions")

	compositions, err := c.repo.ListProcessingCompositions(ctx)
	if err != nil {
		return fmt.Errorf("list processing compositions: %w", err)
	}

	log.Info("poll compositions", "count", len(compositions))

	for _, comp := range compositions {
		if comp.ExternalID == nil || *comp.ExternalID == "" {
			log.Warn("composition has no external ID, skipping", "compositionID", comp.ID)
			continue
		}

		log.Info("checking composition", "compositionID", comp.ID, "externalID", *comp.ExternalID, "videoID", comp.VideoID)

		status, err := c.composer.GetStatus(ctx, *comp.ExternalID)
		if err != nil {
			log.Error("get composition status failed", "compositionID", comp.ID, "externalID", *comp.ExternalID, "videoID", comp.VideoID, "error", err)
			continue
		}
		log.Info("composition status", "compositionID", comp.ID, "externalID", *comp.ExternalID, "status", status.Status)

		switch status.Status {
		case "done":
			comp.Status = blogger.CompositionStatusDone
			comp.ResultURL = &status.ResultURL
			c.notifyDone(ctx, comp.VideoID, status.ResultURL)
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

func (c *PollCompositions) notifyDone(ctx context.Context, videoID, resultURL string) {
	log := logger.FromContext(ctx)

	video, err := c.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		log.Error("poll compositions: get video for notify", "videoID", videoID, "error", err)
		return
	}

	watchers, err := c.repo.ListVideoWatchers(ctx, videoID)
	if err != nil || len(watchers) == 0 {
		return
	}

	evts := make([]any, 0, len(watchers))
	for _, w := range watchers {
		evts = append(evts, &blogger.VideoCompositionDone{
			VideoID:   videoID,
			ChatID:    w.ChatID,
			ResultURL: resultURL,
			SourceURL: video.URL,
			At:        time.Now(),
		})
	}
	c.dispatcher.Dispatch(ctx, evts)
}
