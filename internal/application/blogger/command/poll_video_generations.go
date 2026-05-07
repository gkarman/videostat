package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type PollVideoGenerations struct {
	r       blogger.Repo
	g       application.VideoGenerator
	storage application.Storage
}

func NewPollVideoGenerations(r blogger.Repo, g application.VideoGenerator, storage application.Storage) *PollVideoGenerations {
	return &PollVideoGenerations{r: r, g: g, storage: storage}
}

func (c *PollVideoGenerations) Execute(ctx context.Context) error {
	log := logger.FromContext(ctx).With("component", "PollVideoGenerations")

	pending, err := c.r.ListPendingVideoGenerations(ctx)
	if err != nil {
		return fmt.Errorf("list pending generations: %w", err)
	}

	for _, vg := range pending {
		if vg.Platform != c.g.Platform() {
			continue
		}

		if err := c.poll(ctx, vg); err != nil {
			log.Error("poll failed", "externalID", vg.ExternalID, "error", err)
		}
	}

	return nil
}

func (c *PollVideoGenerations) poll(ctx context.Context, vg *blogger.VideoGeneration) error {
	log := logger.FromContext(ctx).With("externalID", vg.ExternalID)

	status, err := c.g.GetStatus(ctx, vg.ExternalID)
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	switch status.Status {
	case "completed":
		s3URL, err := c.downloadAndUpload(ctx, vg, status.VideoURL)
		if err != nil {
			return fmt.Errorf("download and upload: %w", err)
		}
		vg.Status = blogger.VideoGenerationCompleted
		vg.S3URL = s3URL
		vg.UpdatedAt = time.Now()
		log.Info("generation completed", "s3URL", s3URL)

	case "failed":
		vg.Status = blogger.VideoGenerationFailed
		vg.ErrorMessage = status.ErrorMessage
		vg.UpdatedAt = time.Now()
		log.Warn("generation failed", "reason", status.ErrorMessage)

	case "processing":
		if vg.Status != blogger.VideoGenerationProcessing {
			vg.Status = blogger.VideoGenerationProcessing
			vg.UpdatedAt = time.Now()
		}

	default:
		return nil
	}

	if err := c.r.UpdateVideoGeneration(ctx, vg); err != nil {
		return fmt.Errorf("update video generation: %w", err)
	}

	return nil
}

func (c *PollVideoGenerations) downloadAndUpload(ctx context.Context, vg *blogger.VideoGeneration, videoURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, videoURL, nil)
	if err != nil {
		return "", fmt.Errorf("create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("download error %d: %s", resp.StatusCode, body)
	}

	key := fmt.Sprintf("videos/%s/%s.mp4", vg.VideoID, vg.ExternalID)
	return c.storage.Upload(ctx, key, resp.Body, "video/mp4")
}
