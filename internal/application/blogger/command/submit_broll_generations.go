package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type SubmitBrollGenerations struct {
	repo blogger.Repo
	gen  application.BrollVideoGenerator
}

func NewSubmitBrollGenerations(repo blogger.Repo, gen application.BrollVideoGenerator) *SubmitBrollGenerations {
	return &SubmitBrollGenerations{repo: repo, gen: gen}
}

func (c *SubmitBrollGenerations) Run(ctx context.Context, req reqdto.SubmitBrollGenerations) error {
	log := logger.FromContext(ctx).With("component", "SubmitBrollGenerations", "videoID", req.VideoID)

	segments, err := c.repo.ListPendingBrollSegments(ctx, req.VideoID)
	if err != nil {
		return fmt.Errorf("list pending broll segments: %w", err)
	}

	log.Info("submitting pending broll segments", "count", len(segments))

	for _, s := range segments {
		durationSec := 5
		if (s.EndMS-s.StartMS)/1000 >= 8 {
			durationSec = 10
		}

		externalID, err := c.gen.Submit(ctx, s.BrollPrompt, durationSec)
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "1303") {
				log.Warn("kling rate limit hit, stopping submit — will retry on next cron tick",
					"segmentID", s.ID, "position", s.Position, "error", err)
				break
			}
			if strings.Contains(errStr, "1102") {
				log.Error("kling account balance not enough — top up the account, stopping submit",
					"segmentID", s.ID, "position", s.Position, "error", err)
				break
			}
			log.Error("kling submit failed, marking segment as failed",
				"segmentID", s.ID, "position", s.Position, "error", err)
			s.GenerationStatus = blogger.BrollStatusFailed
			errMsg := errStr
			s.GenerationError = &errMsg
		} else {
			log.Info("kling segment submitted", "segmentID", s.ID, "position", s.Position, "externalID", externalID)
			s.GenerationStatus = blogger.BrollStatusProcessing
			s.GenerationExternalID = &externalID
		}

		if err := c.repo.UpdateBrollSegment(ctx, s); err != nil {
			return fmt.Errorf("update broll segment %s: %w", s.ID, err)
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
