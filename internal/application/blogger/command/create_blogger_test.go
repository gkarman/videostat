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

	repoDictionaryFake := repo_dictionary.NewFake()

	tests := []struct {
		name          string
		req           reqdto.CreateBlogger
		wantErr       error
		expectedEvent any
	}{
		{
			name:          "empty url",
			req:           reqdto.CreateBlogger{URL: "", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
			wantErr:       blogger.ErrUrlInvalid,
			expectedEvent: nil,
		},
		{
			name:          "no exist platform",
			req:           reqdto.CreateBlogger{URL: "https://www.youtube.com/@redactsiya", PlatformName: "some"},
			wantErr:       dictionary.ErrPlatformNotFound,
			expectedEvent: nil,
		},
		{
			name:          "success - valid url and platform",
			req:           reqdto.CreateBlogger{URL: "https://www.youtube.com/@redactsiya", PlatformName: repo_dictionary.PLATFORM_NAME_YOUTUBE},
			wantErr:       nil,
			expectedEvent: blogger.Created{URL: "https://www.youtube.com/@redactsiya"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repoBlogger := repo_blogger.NewFake()
			disp := dispatcher.NewFakeDispatcher()
			cmd := NewCreateBlogger(repoBlogger, repoDictionaryFake, disp)

			resp, err := cmd.Run(ctx, tt.req)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Zero(t, resp)
				require.Empty(t, disp.Events)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, resp)

			if tt.expectedEvent != nil {
				require.Len(t, disp.Events, 1)
				actualEvents := disp.Events[0]
				require.Len(t, actualEvents, 1)

				switch e := tt.expectedEvent.(type) {
				case blogger.Created:
					act := actualEvents[0].(*blogger.Created)
					require.Equal(t, e.URL, act.URL)
					require.NotEmpty(t, act.ID)
				default:
					t.Fatalf("unexpected expectedEvent type: %T", e)
				}
			} else {
				require.Empty(t, disp.Events)
			}
		})
	}
}
