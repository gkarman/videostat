package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	repo_blogger "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/stretchr/testify/require"
)

type fakeAnalyzer struct {
	payload      []byte
	err          error
	providerName string
}

func (f *fakeAnalyzer) Analyze(_ context.Context, _ string) ([]byte, error) {
	return f.payload, f.err
}

func (f *fakeAnalyzer) ProviderName() string {
	return f.providerName
}

func TestAnalyzeVideo_Run(t *testing.T) {
	const videoID = "vid_1"
	const fileURL = "https://storage.example.com/video.webm"

	newProcessingVideo := func() *blogger.Video {
		v := blogger.NewVideo(blogger.CreateVideoDto{
			ID:         videoID,
			ExternalID: "ext_1",
			URL:        "https://youtube.com/v1",
		})
		_ = v.ChangeStatus(blogger.VideoStatusProcessing)
		return v
	}

	tests := []struct {
		name              string
		setupRepo         func(r *repo_blogger.InMemoryRepo)
		analyzer          *fakeAnalyzer
		saveAnalysisErr   error
		wantErr           bool
		wantAnalysesCount int
		wantVideoStatus   *blogger.VideoStatus
		wantEvents        int
	}{
		{
			name: "успех — анализ сохранён, событие VideoAnalyzeDone выброшено",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				_ = r.SaveVideo(context.Background(), newProcessingVideo())
			},
			analyzer:          &fakeAnalyzer{payload: []byte(`{"text":"hello"}`), providerName: "fake"},
			wantErr:           false,
			wantAnalysesCount: 1,
			wantEvents:        1,
		},
		{
			name: "видео не найдено — возвращает ошибку",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				// репо пустое
			},
			analyzer:          &fakeAnalyzer{payload: []byte(`{}`), providerName: "fake"},
			wantErr:           true,
			wantAnalysesCount: 0,
		},
		{
			name: "ошибка анализатора — видео переходит в статус failed",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				_ = r.SaveVideo(context.Background(), newProcessingVideo())
			},
			analyzer:          &fakeAnalyzer{err: errors.New("assemblyai timeout"), providerName: "fake"},
			wantErr:           true,
			wantAnalysesCount: 0,
			wantVideoStatus:   statusPtr(blogger.VideoStatusFailed),
		},
		{
			name: "ошибка сохранения анализа — возвращает ошибку",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				_ = r.SaveVideo(context.Background(), newProcessingVideo())
			},
			analyzer:          &fakeAnalyzer{payload: []byte(`{}`), providerName: "fake"},
			saveAnalysisErr:   errors.New("db error"),
			wantErr:           true,
			wantAnalysesCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repo_blogger.NewInMemory()
			repo.SaveVideoAnalysisErr = tt.saveAnalysisErr
			tt.setupRepo(repo)

			disp := dispatcher.NewFakeDispatcher()
			cmd := command.NewAnalyzeVideo(repo, tt.analyzer, disp)

			err := cmd.Run(context.Background(), reqdto.AnalyzeVideo{
				VideoID: videoID,
				FileURL: fileURL,
			})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			analyses := repo.GetSavedAnalyses()
			require.Len(t, analyses, tt.wantAnalysesCount)

			if tt.wantAnalysesCount > 0 {
				a := analyses[0]
				require.Equal(t, videoID, a.VideoID)
				require.Equal(t, tt.analyzer.providerName, a.Provider)
				require.Equal(t, tt.analyzer.payload, a.RawPayload)
			}

			totalEvents := 0
			for _, batch := range disp.Events {
				totalEvents += len(batch)
			}
			require.Equal(t, tt.wantEvents, totalEvents)

			if tt.wantVideoStatus != nil {
				v, err := repo.GetVideoByID(context.Background(), videoID)
				require.NoError(t, err)
				require.Equal(t, *tt.wantVideoStatus, v.Status)
			}
		})
	}
}

func statusPtr(s blogger.VideoStatus) *blogger.VideoStatus {
	return &s
}
