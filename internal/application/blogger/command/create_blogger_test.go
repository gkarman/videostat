package command

import (
	"context"
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/domain/dictionary"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	repo_blogger "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	repo_dictionary "github.com/gkarman/demo/internal/infrastructure/repository/dictionary"
	"github.com/stretchr/testify/require"
)

func TestCreateBlogger(t *testing.T) {
	ctx := context.Background()
	repoDictionaryFake := repo_dictionary.NewFake()

	tests := []struct {
		name          string
		reqs          []reqdto.CreateBlogger // несколько запросов подряд, чтобы проверить дубликаты
		wantErrs      []error
		expectedEvent []blogger.Created
	}{
		{
			name: "empty url",
			reqs: []reqdto.CreateBlogger{
				{URL: "", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
			},
			wantErrs:      []error{blogger.ErrUrlInvalid},
			expectedEvent: nil,
		},
		{
			name: "no exist platform",
			reqs: []reqdto.CreateBlogger{
				{URL: "https://www.youtube.com/@redactsiya", PlatformName: "some"},
			},
			wantErrs:      []error{dictionary.ErrPlatformNotFound},
			expectedEvent: nil,
		},
		{
			name: "success - valid url and platform",
			reqs: []reqdto.CreateBlogger{
				{URL: "https://www.youtube.com/@redactsiya", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
			},
			wantErrs:      []error{nil},
			expectedEvent: []blogger.Created{{URL: "https://www.youtube.com/@redactsiya"}},
		},
		{
			name: "duplicate url",
			reqs: []reqdto.CreateBlogger{
				{URL: "https://www.youtube.com/@redactsiya", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
				{URL: "https://www.youtube.com/@redactsiya", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
			},
			wantErrs:      []error{nil, blogger.ErrUrlExist},
			expectedEvent: []blogger.Created{{URL: "https://www.youtube.com/@redactsiya"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoBlogger := repo_blogger.NewInMemory()
			disp := dispatcher.NewFakeDispatcher()
			cmd := NewCreateBlogger(repoBlogger, repoDictionaryFake, disp)

			for i, req := range tt.reqs {
				resp, err := cmd.Run(ctx, req)

				if tt.wantErrs[i] != nil {
					require.ErrorIs(t, err, tt.wantErrs[i])
					require.Zero(t, resp)
					continue
				}

				require.NoError(t, err)
				require.NotZero(t, resp)
			}

			if len(tt.expectedEvent) > 0 {
				require.Len(t, disp.Events, 1)
				actualEvents := disp.Events[0]
				require.Len(t, actualEvents, len(tt.expectedEvent))

				for j, e := range tt.expectedEvent {
					act := actualEvents[j].(*blogger.Created)
					require.Equal(t, e.URL, act.URL)
					require.NotEmpty(t, act.ID)
				}
			} else {
				require.Empty(t, disp.Events)
			}
		})
	}
}