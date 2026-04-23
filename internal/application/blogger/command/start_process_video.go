package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/application/blogger/command/respdto"
	"github.com/gkarman/demo/internal/domain/blogger"
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

	oldStatus := v.Status
    err = v.StartProcessing()
	if err != nil {
		return respdto.StartProcessVideo{}, fmt.Errorf("start processing domain logic: %w", err)
	}

	err = c.repo.UpdateVideoStatus(ctx, v.ID, oldStatus, v.Status)
	if err != nil {
		return respdto.StartProcessVideo{}, fmt.Errorf("update video status in db: %w", err)	}

	c.disp.Dispatch(ctx, v.PullEvents())
	return respdto.StartProcessVideo{
			Message: "Начался процесс сбора данных...",	},
		nil
}
