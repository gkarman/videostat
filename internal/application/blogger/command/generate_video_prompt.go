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

type GenerateVideoPrompt struct {
	r blogger.Repo
	g application.VideoPromptGenerator
}

func NewGenerateVideoPrompt(r blogger.Repo, g application.VideoPromptGenerator) *GenerateVideoPrompt {
	return &GenerateVideoPrompt{r: r, g: g}
}

func (c *GenerateVideoPrompt) Run(ctx context.Context, req reqdto.GenerateVideoPrompt) error {
	log := logger.FromContext(ctx).With("component", "GenerateVideoPrompt", "videoID", req.VideoID)

	va, err := c.r.GetVideoAnalysisByVideoID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video analysis: %w", err)
	}

	prompt, err := c.g.Generate(ctx, va.RawPayload)
	if err != nil {
		return fmt.Errorf("generate prompt: %w", err)
	}

	vp := &blogger.VideoPrompt{
		ID:          uuid.NewString(),
		VideoID:     req.VideoID,
		LLMProvider: c.g.ProviderName(),
		LLMModel:    c.g.ModelName(),
		Prompt:      prompt,
		CreatedAt:   time.Now(),
	}

	if err = c.r.SaveVideoPrompt(ctx, vp); err != nil {
		return fmt.Errorf("save video prompt: %w", err)
	}

	v, err := c.r.GetVideoByID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video: %w", err)
	}

	if err = v.MarkReady(); err != nil {
		return fmt.Errorf("mark ready: %w", err)
	}

	if err = c.r.UpdateVideoState(ctx, v); err != nil {
		return fmt.Errorf("update video state: %w", err)
	}

	log.Info("video prompt generated", "provider", vp.LLMProvider, "model", vp.LLMModel)
	return nil
}
