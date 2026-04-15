package query

import "github.com/gkarman/demo/internal/application/blogger/query/view"

type VideoMetricsEnricher interface {
	Enrich(videos []*view.Video)
}