package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
)

type GenerateBrollSegments struct {
	repo   blogger.Repo
	gen    application.BrollGenerator
	submit *SubmitBrollGenerations
}

func NewGenerateBrollSegments(repo blogger.Repo, gen application.BrollGenerator, submit *SubmitBrollGenerations) *GenerateBrollSegments {
	return &GenerateBrollSegments{repo: repo, gen: gen, submit: submit}
}

func (c *GenerateBrollSegments) Run(ctx context.Context, req reqdto.GenerateBrollSegments) error {
	va, err := c.repo.GetVideoAnalysisByVideoID(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("get video analysis: %w", err)
	}

	segments, err := c.gen.GenerateSegments(ctx, va.RawPayload)
	if err != nil {
		return fmt.Errorf("generate broll segments: %w", err)
	}

	if err := c.repo.SaveBrollSegments(ctx, req.VideoID, segments); err != nil {
		return fmt.Errorf("save broll segments: %w", err)
	}

	if err := c.submit.Run(ctx, reqdto.SubmitBrollGenerations{VideoID: req.VideoID}); err != nil {
		return fmt.Errorf("submit broll generations: %w", err)
	}

	return nil
}
