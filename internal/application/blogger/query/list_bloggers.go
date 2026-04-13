package query

import (
	"context"

	"github.com/gkarman/demo/internal/application/blogger"
	"github.com/gkarman/demo/internal/application/blogger/query/respdto"
	"github.com/gkarman/demo/internal/application/blogger/query/view"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type ListBloggers struct {
	repo blogger.ReadRepo
}

func NewListBloggers(repo blogger.ReadRepo) *ListBloggers {
	return &ListBloggers{repo: repo}
}

func (q *ListBloggers) Run(ctx context.Context) (*respdto.ListBloggers, error) {
	log := logger.FromContext(ctx).With("component", "ListBloggers")
	log.Debug("start query list blogger")

	rows, err := q.repo.ListBloggers(ctx)
	if err != nil {
		log.Error("list rows failed", "error", err)
		return nil, err
	}

	items := make([]*view.Blogger, 0, len(rows))
	for _, r := range rows {
		items = append(items, &view.Blogger{
			ID:       r.ID,
			URL:      r.URL,
			Platform: r.Platform,
		})
	}
	log.Debug("find bloggers", "items", items)

	return &respdto.ListBloggers{
		Items: items,
	}, nil
}
