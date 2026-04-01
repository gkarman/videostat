package video

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application/video/requestdto"
	"github.com/gkarman/demo/internal/domain/video"
)

type CreateAccount struct {
	repo video.Repository
}

func NewCreateAccountService(repo video.Repository) *CreateAccount {
	return &CreateAccount{
		repo: repo,
	}
}

func (s *CreateAccount) Execute(ctx context.Context, req requestdto.CreateAccount) error {
	exists, err := s.repo.ExistsByPlatformAndExternalID(ctx, req.PlatformID, req.ExternalID)
	if err != nil {
		return fmt.Errorf("repository ExistsByPlatformAndExternalID: %w", err)
	}
	if exists {
		return video.ErrAccountAlreadyExists
	}

	snap := &video.AccountSnapshot{
		Account: &video.Account{
			PlatformID: req.PlatformID,
			ExternalID: req.ExternalID,
			Title:      req.Title,
			URL:        req.URL,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		AccountStats: &video.AccountStats{
			CollectedAt: time.Now(),
		},
		Contents: []*video.ContentSnapshot{},
	}

	err = s.repo.SaveSnapshot(ctx, snap)
	if err != nil {
		return fmt.Errorf("repository SaveSnapshot: %w", err)
	}

	return nil
}
