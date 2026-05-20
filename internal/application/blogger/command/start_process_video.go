package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/application/blogger/command/respdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/google/uuid"
)

type StartProcessVideo struct {
	repo blogger.Repo
	disp application.Dispatcher
}

func NewStartProcessVideo(r blogger.Repo, d application.Dispatcher) *StartProcessVideo {
	return &StartProcessVideo{
		repo: r,
		disp: d,
	}
}

func (c *StartProcessVideo) Run(ctx context.Context, req reqdto.StartProcessVideo) (respdto.StartProcessVideo, error) {
	v, err := c.repo.GetVideoByUrl(ctx, req.URL)
	if err != nil {
		return respdto.StartProcessVideo{}, fmt.Errorf("get video by url: %w", err)
	}

	if req.ChatID != 0 {
		w := &blogger.VideoWatcher{
			ID:        uuid.NewString(),
			VideoID:   v.ID,
			ChatID:    req.ChatID,
			CreatedAt: time.Now(),
		}
		if err = c.repo.SaveVideoWatcher(ctx, w); err != nil {
			return respdto.StartProcessVideo{}, fmt.Errorf("save video watcher: %w", err)
		}
	}

	_, err = c.repo.GetVideoAnalysisByVideoID(ctx, v.ID)
	if err == nil {
		v.MarkForRegeneration()
		if err = c.repo.UpdateVideoState(ctx, v); err != nil {
			return respdto.StartProcessVideo{}, fmt.Errorf("update video state: %w", err)
		}
		v.AnalyzeDone()
		c.disp.Dispatch(ctx, v.PullEvents())
		return respdto.StartProcessVideo{Message: "Генерация запущена..."}, nil
	}

	oldStatus := v.Status
	if err = v.MarkStartProcessing(); err != nil {
		return respdto.StartProcessVideo{}, fmt.Errorf("start processing domain logic: %w", err)
	}

	if err = c.repo.UpdateVideoStatus(ctx, v.ID, oldStatus, v.Status); err != nil {
		return respdto.StartProcessVideo{}, fmt.Errorf("update video status in db: %w", err)
	}

	c.disp.Dispatch(ctx, v.PullEvents())
	return respdto.StartProcessVideo{Message: "Начался процесс сбора данных..."}, nil
}
