package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/reqdto"
	"github.com/gkarman/demo/internal/application/blogger/respdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/domain/dictionary"
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
	platform, err := c.repoDictionary.GetPlatformByName(ctx, req.PlatformName)
	if err != nil {
		return respdto.CreateBlogger{}, fmt.Errorf("get platform by name in createBlogger: %w", err)
	}

	if platform == nil {
		return respdto.CreateBlogger{}, dictionary.ErrPlatformNotFound
	}

	r := blogger.CreateBloggerDto{
		ID:         uuid.NewString(),
		PlatformID: platform.ID,
		URL:        req.URL,
	}

	existBloggerUrl, err := c.repoBlogger.ExistByUrl(ctx, r.URL)
	if err != nil {
		return respdto.CreateBlogger{}, fmt.Errorf("check exist in db: %w", err)
	}

	if existBloggerUrl {
		return respdto.CreateBlogger{}, blogger.ErrUrlExist
	}

	b, err := blogger.Create(r)
	if err != nil {
		return respdto.CreateBlogger{}, fmt.Errorf("create blogger: %w", err)
	}

	if err := c.repoBlogger.Save(ctx, b); err != nil {
		return respdto.CreateBlogger{}, fmt.Errorf("save blogger in db: %w", err)
	}

	events := b.PullEvents()
	c.disp.Dispatch(ctx, events)

	resp := respdto.CreateBlogger{
		ID: b.ID,
	}

	return resp, nil
}
