package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/llm"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type GenerateVideoPrompt struct {
	r blogger.Repo
	g application.VideoPromptGenerator
	d application.Dispatcher
}

func NewGenerateVideoPrompt(r blogger.Repo, g application.VideoPromptGenerator, d application.Dispatcher) *GenerateVideoPrompt {
	return &GenerateVideoPrompt{r: r, g: g, d: d}
}

func (c *GenerateVideoPrompt) Run(ctx context.Context, req reqdto.GenerateVideoPrompt) error {
	log := logger.FromContext(ctx).With("component", "GenerateVideoPrompt", "videoID", req.VideoID)

	v, err := c.r.GetVideoByID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video: %w", err)
	}

	va, err := c.r.GetVideoAnalysisByVideoID(ctx, req.VideoID)
	if err != nil {
		return c.failVideo(ctx, v, blogger.ErrorStageInsights, fmt.Errorf("get video analysis: %w", err))
	}

	raw, err := c.g.Generate(ctx, va.RawPayload)
	if err != nil {
		return c.failVideo(ctx, v, blogger.ErrorStageInsights, fmt.Errorf("generate prompt: %w", err))
	}

	content, err := llm.ParseGeneratedContent(raw)
	if err != nil {
		return c.failVideo(ctx, v, blogger.ErrorStageInsights, fmt.Errorf("parse llm response: %w", err))
	}

	vp := &blogger.VideoPrompt{
		ID:          uuid.NewString(),
		VideoID:     req.VideoID,
		LLMProvider: c.g.ProviderName(),
		LLMModel:    c.g.ModelName(),
		Brief:       string(content.Brief),
		Script:      content.Script,
		CreatedAt:   time.Now(),
	}

	if err = c.r.SaveVideoPrompt(ctx, vp); err != nil {
		return c.failVideo(ctx, v, blogger.ErrorStageInsights, fmt.Errorf("save video prompt: %w", err))
	}

	if err = v.MarkReady(); err != nil {
		return fmt.Errorf("mark ready: %w", err)
	}

	v.PromptGenerated()

	if err = c.r.UpdateVideoState(ctx, v); err != nil {
		return fmt.Errorf("update video state: %w", err)
	}

	c.d.Dispatch(ctx, v.PullEvents())

	log.Info("video prompt generated", "provider", vp.LLMProvider, "model", vp.LLMModel)
	return nil
}

func (c *GenerateVideoPrompt) failVideo(ctx context.Context, v *blogger.Video, stage blogger.VideoErrorStage, cause error) error {
	if err := v.MarkFailProcessing(stage, cause); err != nil {
		return fmt.Errorf("%w (also failed to mark video as failed: %v)", cause, err)
	}
	if err := c.r.UpdateVideoState(ctx, v); err != nil {
		return fmt.Errorf("%w (also failed to persist video failure: %v)", cause, err)
	}
	return cause
}
