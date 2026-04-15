package analytics

import (
	"github.com/gkarman/demo/internal/application/blogger/query"
	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

type VideoMetricsEnricher struct {
	enrichers []query.VideoEnricher
}

func NewVideoMetricsEnricher(enrichers ...query.VideoEnricher) *VideoMetricsEnricher {
	return &VideoMetricsEnricher{
		enrichers: enrichers,
	}
}

func (e *VideoMetricsEnricher) Enrich(videos []*view.Video) {
	for _, enricher := range e.enrichers {
		enricher.Enrich(videos)
	}
}
