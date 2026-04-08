package command

import (
	"context"
	"fmt"
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher"
	"github.com/stretchr/testify/require"
)

func TestFetchBloggerVideos(t *testing.T) {
	ctx := context.Background()

	baseVideos := []*blogger.Video{
		{
			ID:         "v1",
			ExternalID: "ext1",
			URL:        "url1",
			Title:      "Video 1",
		},
		{
			ID:         "v2",
			ExternalID: "ext2",
			URL:        "url2",
			Title:      "Video 2",
		},
	}

	tests := []struct {
		name           string
		bloggerExists  bool
		searchErr      error
		saveVideoErr   map[string]error
		wantErr        bool
		wantSavedCount int
	}{
		{
			name:           "blogger not found",
			bloggerExists:  false,
			wantErr:        true,
			wantSavedCount: 0,
		},
		{
			name:           "searcher error",
			bloggerExists:  true,
			searchErr:      fmt.Errorf("apify error"),
			wantErr:        true,
			wantSavedCount: 0,
		},
		{
			name:          "save video partially fails but command succeeds",
			bloggerExists: true,
			saveVideoErr: map[string]error{
				"ext1": fmt.Errorf("db error"),
			},
			wantErr:        false,
			wantSavedCount: 1,
		},
		{
			name:           "success - all videos saved",
			bloggerExists:  true,
			wantErr:        false,
			wantSavedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repo.NewInMemory()

			if tt.bloggerExists {
				b := &blogger.Blogger{
					ID:         "b1",
					PlatformID: 1,
					URL:        "url",
				}
				require.NoError(t, repo.Save(ctx, b))
			}

			searcher := videosearcher.NewInMemoryVideoSearcher(baseVideos, tt.searchErr)

			if tt.saveVideoErr != nil {
				repo.SaveVideoErrFor = tt.saveVideoErr
			}

			cmd := NewFetchBloggerVideos(repo, searcher)

			err := cmd.Execute(ctx, reqdto.FetchBloggerVideos{
				BloggerID: "b1",
			})

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			saved, err := repo.ListVideosByBlogger(ctx, "b1")
			require.NoError(t, err)
			require.Len(t, saved, tt.wantSavedCount)
		})
	}
}