package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	blogger_repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	videosrc "github.com/gkarman/demo/internal/infrastructure/videosourcesearcher"
)

func TestFetchVideoSources(t *testing.T) {
	tests := []struct {
		name          string
		searchErr     error
		searchResult  string
		wantStatus    blogger.VideoStatus
		wantErr       bool
		wantEventSize int
	}{
		{
			name:          "success",
			searchResult:  "file-url",
			wantStatus:    blogger.VideoStatusProcessing,
			wantErr:       false,
			wantEventSize: 1,
		},
		{
			name:       "search error",
			searchErr:  assertError(),
			wantStatus: blogger.VideoStatusFailed,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := blogger_repo.NewInMemory()

			video := &blogger.Video{
				ID:         "1",
				URL:        "url-1",
				ExternalID: "ext-1",
				Status:     blogger.VideoStatusProcessing,
			}

			_ = repo.SaveVideo(context.Background(), video)

			searcher := videosrc.NewInMemorySourceSearcher()

			if tt.searchErr != nil {
				searcher.Errs["1"] = tt.searchErr
			} else {
				searcher.Data["1"] = tt.searchResult
			}

			disp := dispatcher.NewFakeDispatcher()

			cmd := FetchVideoSources{
				repo:     repo,
				searcher: searcher,
				disp:     disp,
			}

			err := cmd.Execute(context.Background(), reqdto.FetchVideoSources{
				VideoID:  "1",
				VideoURL: "url-1",
			})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			v, err := repo.GetVideoByID(context.Background(), "1")
			require.NoError(t, err)

			require.Equal(t, tt.wantStatus, v.Status)

			if tt.wantEventSize > 0 {
				require.Len(t, disp.Events, 1)
				require.Len(t, disp.Events[0], tt.wantEventSize)
			}
		})
	}
}

func assertError() error {
	return &testErr{}
}

type testErr struct{}

func (e *testErr) Error() string { return "boom" }
