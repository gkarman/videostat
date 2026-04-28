package command

import (
	"context"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command/reqdto"
	"github.com/gkarman/demo/internal/domain/blogger"
)

type AnalyzeVideo struct {
	r blogger.Repo
	d application.Dispatcher
}

func NewAnalyzeVideo(r blogger.Repo, d application.Dispatcher) *AnalyzeVideo {
	return &AnalyzeVideo{
		r: r,
		d: d,
	}
}

func (c *AnalyzeVideo) Run(ctx context.Context, req reqdto.AnalyzeVideo) error {
	return nil
}
