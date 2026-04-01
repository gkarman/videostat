package video

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/video/requestdto"
	"github.com/gkarman/demo/internal/domain/video"
)

type CreateAccountSnapshot struct {
	repo  video.Repository
	apify application.ApifyClient
}

func NewCreateAccountSnapshotService(repo video.Repository, apify application.ApifyClient) *CreateAccountSnapshot {
	return &CreateAccountSnapshot{
		repo:  repo,
		apify: apify,
	}
}

func (s *CreateAccountSnapshot) Execute(ctx context.Context, req requestdto.CreateAccountSnapshot) error {
	snap, err := s.apify.FetchAccountSnapshot(ctx, req.AccountURL)
	if err != nil {
		return fmt.Errorf("fetch snapshot: %w", err)
	}

	return s.repo.SaveSnapshot(ctx, snap)
}