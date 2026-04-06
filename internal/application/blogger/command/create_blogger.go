package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/reqdto"
	"github.com/gkarman/demo/internal/application/blogger/respdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/domain/dictionary"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type CreateBlogger struct {
	repoBlogger    blogger.Repo
	repoDictionary dictionary.Repo
	disp           application.Dispatcher
}

func NewCreateBlogger(
	repoBlogger blogger.Repo,
	repoDictionary dictionary.Repo,
	disp application.Dispatcher,
) *CreateBlogger {
	return &CreateBlogger{
		repoBlogger:    repoBlogger,
		repoDictionary: repoDictionary,
		disp:           disp,
	}
}

func (c *CreateBlogger) Run(ctx context.Context, req reqdto.CreateBlogger) (respdto.CreateBlogger, error) {
	log := logger.FromContext(ctx).With(
		"component", "CreateBlogger",
		"url", req.URL,
		"platform", req.PlatformName,
	)

	log.Debug("start")

	platform, err := c.repoDictionary.GetPlatformByName(ctx, req.PlatformName)
	if err != nil {
		log.Error("get platform by name failed", "err", err)
		return respdto.CreateBlogger{}, fmt.Errorf("get platform by name: %w", err)
	}
	if platform == nil {
		log.Warn("platform not found")
		return respdto.CreateBlogger{}, dictionary.ErrPlatformNotFound
	}

	exist, err := c.repoBlogger.ExistByUrl(ctx, req.URL)
	if err != nil {
		log.Error("exist by url failed", "err", err)
		return respdto.CreateBlogger{}, fmt.Errorf("exist by url: %w", err)
	}
	if exist {
		log.Warn("url already exists")
		return respdto.CreateBlogger{}, blogger.ErrUrlExist
	}

	b, err := blogger.Create(blogger.CreateBloggerDto{
		ID:         uuid.NewString(),
		PlatformID: platform.ID,
		URL:        req.URL,
	})
	if err != nil {
		log.Error("create blogger failed", "err", err)
		return respdto.CreateBlogger{}, fmt.Errorf("create blogger: %w", err)
	}

	if err := c.repoBlogger.Save(ctx, b); err != nil {
		log.Error("save blogger failed", "blogger_id", b.ID, "err", err)
		return respdto.CreateBlogger{}, fmt.Errorf("save blogger: %w", err)
	}

	c.disp.Dispatch(ctx, b.PullEvents())

	log.Info("success", "blogger_id", b.ID)

	return respdto.CreateBlogger{ID: b.ID}, nil
}