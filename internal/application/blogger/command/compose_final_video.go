package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/google/uuid"
)

var ErrAvatarNotReady = errors.New("avatar video not ready yet")

type ComposeFinalVideo struct {
	repo     blogger.Repo
	composer application.VideoComposer
}

func NewComposeFinalVideo(repo blogger.Repo, composer application.VideoComposer) *ComposeFinalVideo {
	return &ComposeFinalVideo{repo: repo, composer: composer}
}

func (c *ComposeFinalVideo) Run(ctx context.Context, req reqdto.ComposeFinalVideo) error {
	gen, err := c.repo.GetVideoGenerationByVideoID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video generation: %w", err)
	}
	if gen.S3URL == "" {
		return ErrAvatarNotReady
	}

	segments, err := c.repo.ListBrollSegmentsByVideoID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("list broll segments: %w", err)
	}

	clips := make([]application.CompositionClip, 0, len(segments))
	var cursor float64
	for _, s := range segments {
		if s.GenerationStatus != blogger.BrollStatusDone || s.BrollURL == nil {
			continue
		}
		// trim Kling clip to actual segment duration for fast-paced cuts
		clipSec := float64(s.EndMS-s.StartMS) / 1000
		if clipSec < 1.0 {
			clipSec = 1.0
		}
		clips = append(clips, application.CompositionClip{
			URL:      *s.BrollURL,
			StartSec: cursor,
			LenSec:   clipSec,
		})
		cursor += clipSec
	}
	totalSec := cursor

	if len(clips) == 0 {
		return fmt.Errorf("no done broll segments to compose")
	}

	externalID, err := c.composer.Compose(ctx, application.CompositionRequest{
		AvatarURL:  gen.S3URL,
		BrollClips: clips,
		TotalSec:   totalSec,
	})
	if err != nil {
		return fmt.Errorf("submit composition: %w", err)
	}

	composition := &blogger.VideoComposition{
		ID:         uuid.NewString(),
		VideoID:    req.VideoID,
		Status:     blogger.CompositionStatusProcessing,
		ExternalID: &externalID,
		CreatedAt:  time.Now(),
	}

	return c.repo.SaveVideoComposition(ctx, composition)
}
