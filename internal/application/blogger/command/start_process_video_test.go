package command_test

import (
	"context"
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	repo_blogger "github.com/gkarman/demo/internal/infrastructure/repository/blogger"

)

func TestStartProcessVideo_Run(t *testing.T) {
	type testCase struct {
		name         string
		inputUrl     string
		setupRepo    func(r *repo_blogger.InMemoryRepo)
		wantErr      bool
		expectedStatus blogger.VideoStatus
		expectedEvents int
	}

	testCases := []testCase{
		{
			name:     "Успешный старт: статус меняется на Processing и событие отправлено",
			inputUrl: "https://youtube.com/v1",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				v := blogger.NewVideo(blogger.CreateVideoDto{
					ID:         "vid_1",
					ExternalID: "ext_1",
					URL:        "https://youtube.com/v1",
				})
				_ = r.SaveVideo(context.Background(), v)
			},
			wantErr:        false,
			expectedStatus: blogger.VideoStatusProcessing,
			expectedEvents: 1,
		},
		{
			name:     "Ошибка: видео не найдено в базе",
			inputUrl: "https://unknown.com",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				// пусто
			},
			wantErr:        true,
			expectedEvents: 0,
		},
		{
			name:     "Ошибка: видео уже готово (невалидный переход)",
			inputUrl: "https://youtube.com/ready",
			setupRepo: func(r *repo_blogger.InMemoryRepo) {
				v := blogger.NewVideo(blogger.CreateVideoDto{
					ID:         "vid_2",
					ExternalID: "ext_2",
					URL:        "https://youtube.com/ready",
				})

				// 1. Сначала переводим в промежуточный статус (Created -> Processing)
				if err := v.ChangeStatus(blogger.VideoStatusProcessing); err != nil {
					t.Fatalf("setup failed: could not change status to Processing: %v", err)
				}

				// 2. Теперь в финальный для этого теста статус (Processing -> Ready)
				if err := v.ChangeStatus(blogger.VideoStatusReady); err != nil {
					t.Fatalf("setup failed: could not change status to Ready: %v", err)
				}

				// 3. Сохраняем уже подготовленный объект
				if err := r.SaveVideo(context.Background(), v); err != nil {
					t.Fatalf("setup failed: could not save video: %v", err)
				}
			},
			wantErr:        true,
			expectedEvents: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := repo_blogger.NewInMemory()
			disp := dispatcher.NewFakeDispatcher()
			tc.setupRepo(repo)

			cmd := command.NewStartProcessVideo(repo, disp)

			// Act
			_, err := cmd.Run(context.Background(), reqdto.StartProcessVideo{URL: tc.inputUrl})

			// Assert
			if (err != nil) != tc.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				// 1. Проверяем статус в БД через репозиторий
				v, err := repo.GetVideoByUrl(context.Background(), tc.inputUrl)
				if err != nil {
					t.Fatalf("could not find video after update: %v", err)
				}
				if v.Status != tc.expectedStatus {
					t.Errorf("expected status %v, got %v", tc.expectedStatus, v.Status)
				}

				// 2. Проверяем события в твоем FakeDispatcher
				// Считаем общее кол-во событий во всех вызовах Dispatch
				totalEvents := 0
				for _, batch := range disp.Events {
					totalEvents += len(batch)
				}

				if totalEvents != tc.expectedEvents {
					t.Errorf("expected %d events total, got %d", tc.expectedEvents, totalEvents)
				}
			}
		})
	}
}