package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
)

type SubmitBrollGenerations struct {
	repo blogger.Repo
	gen  application.BrollVideoGenerator
}

func NewSubmitBrollGenerations(repo blogger.Repo, gen application.BrollVideoGenerator) *SubmitBrollGenerations {
	return &SubmitBrollGenerations{repo: repo, gen: gen}
}

func (c *SubmitBrollGenerations) Run(ctx context.Context, req reqdto.SubmitBrollGenerations) error {
	segments, err := c.repo.ListPendingBrollSegments(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("list pending broll segments: %w", err)
	}

	for _, s := range segments {
		durationSec := 5
		if (s.EndMS-s.StartMS)/1000 >= 8 {
			durationSec = 10
		}

		externalID, err := c.gen.Submit(ctx, s.BrollPrompt, durationSec)
		if err != nil {
			s.GenerationStatus = blogger.BrollStatusFailed
			errMsg := err.Error()
			s.GenerationError = &errMsg
		} else {
			s.GenerationStatus = blogger.BrollStatusProcessing
			s.GenerationExternalID = &externalID
		}

		if err := c.repo.UpdateBrollSegment(ctx, s); err != nil {
			return fmt.Errorf("update broll segment %s: %w", s.ID, err)
		}
	}

	return nil
}
