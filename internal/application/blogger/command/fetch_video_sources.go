package command

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type FetchVideoSources struct {
	repo     blogger.Repo
	searcher application.FileSourceSearcher
	disp     application.Dispatcher
}

func NewFetchVideoSources(repo blogger.Repo, searcher application.FileSourceSearcher, disp application.Dispatcher) *FetchVideoSources {
	return &FetchVideoSources{
		repo:     repo,
		searcher: searcher,
		disp:     disp,
	}
}

func (c *FetchVideoSources) Execute(ctx context.Context, req reqdto.FetchVideoSources) error {
	log := logger.FromContext(ctx).With(
		"component", "FetchVideoSources",
		"videoID", req.VideoID,
	)
	log.Debug("StartFetching")

	v, err := c.repo.GetVideoByUrl(ctx, req.VideoURL)
	if err != nil {
		return fmt.Errorf("get video by id: %w", err)
	}

	fileUrl, err := c.searcher.SearchUrl(ctx, v)
	if err != nil {
		log.Debug("ошибка получения fileUrl из стороннего сервиса", "err",  err)
		err = v.MarkFailProcessing(blogger.ErrorStageFileFetch, err)
		if err != nil {
			return fmt.Errorf("mark fail processing: %w", err)
		}

		err = c.repo.UpdateVideoState(ctx,v)
		if err != nil {
			return fmt.Errorf("update in db state : %w", err)
		}

		return fmt.Errorf("search url: %w", err)
	}
	log.Debug("find url", "url", fileUrl)

	v.SourceFound(fileUrl)
	c.disp.Dispatch(ctx, v.PullEvents())

	return nil
}
