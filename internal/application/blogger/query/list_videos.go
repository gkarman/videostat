package query

import (
	"context"

	"github.com/gkarman/demo/internal/application/blogger"
	"github.com/gkarman/demo/internal/application/blogger/query/respdto"
	"github.com/gkarman/demo/internal/application/blogger/query/view"
	"github.com/gkarman/demo/internal/infrastructure/logger"
)

type ListVideos struct {
	repo     blogger.ReadRepo
	enricher VideoEnricher
}

func NewListVideos(repo blogger.ReadRepo, enricher VideoEnricher) *ListVideos {
	return &ListVideos{
		repo:     repo,
		enricher: enricher,
	}
}

func (q *ListVideos) Run(ctx context.Context) (*respdto.ListVideos, error) {
	log := logger.FromContext(ctx).With("component", "ListVideos")
	log.Debug("start query list videos")

	rows, err := q.repo.ListVideos(ctx)
	if err != nil {
		log.Error("list rows failed", "error", err)
		return nil, err
	}

	items := make([]*view.Video, 0, len(rows))
	for _, r := range rows {
		items = append(items, &view.Video{
			ID:          r.ID,
			Platform:    r.Platform,
			BloggerURL:  r.BloggerURL,
			URL:         r.URL,
			Title:       r.Title,
			Views:       r.Views,
			Likes:       r.Likes,
			Comments:    r.Comments,
			PublishedAt: r.PublishedAt,
			CreatedAt:   r.CreatedAt,
		})
	}
	q.enricher.Enrich(items)
	log.Debug("find bloggers", "items", items)

	return &respdto.ListVideos{
		Items: items,
	}, nil
}
