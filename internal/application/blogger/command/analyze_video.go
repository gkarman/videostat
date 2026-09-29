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

type AnalyzeVideo struct {
	r blogger.Repo
	a application.VideoAnalyzer
	d application.Dispatcher
}

func NewAnalyzeVideo(r blogger.Repo, a application.VideoAnalyzer, d application.Dispatcher) *AnalyzeVideo {
	return &AnalyzeVideo{
		r: r,
		a: a,
		d: d,
	}
}

func (c *AnalyzeVideo) Run(ctx context.Context, req reqdto.AnalyzeVideo) error {
	ctx = logger.WithField(ctx, logger.KeyVideoID, req.VideoID)
	log := logger.FromContext(ctx).With("component", "AnalyzeVideo")

	v, err := c.r.GetVideoByID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video: %w", err)
	}

	raw, err := c.a.Analyze(ctx, req.FileURL)
	if err != nil {
		log.Error("analyze failed", "error", err)
		ferr := v.MarkFailProcessing(blogger.ErrorStageAnalysis, err)
		if ferr != nil {
			return fmt.Errorf("mark fail: %w", ferr)
		}
		if ferr = c.r.UpdateVideoState(ctx, v); ferr != nil {
			return fmt.Errorf("update video state: %w", ferr)
		}
		return fmt.Errorf("analyze: %w", err)
	}

	va := &blogger.VideoAnalysis{
		ID:         uuid.NewString(),
		VideoID:    v.ID,
		Provider:   c.a.ProviderName(),
		RawPayload: raw,
		CreatedAt:  time.Now(),
	}

	if err = c.r.SaveVideoAnalysis(ctx, va); err != nil {
		return fmt.Errorf("save video analysis: %w", err)
	}

	v.AnalyzeDone()
	c.d.Dispatch(ctx, v.PullEvents())

	log.Info("video analysis saved", "provider", va.Provider)
	return nil
}
