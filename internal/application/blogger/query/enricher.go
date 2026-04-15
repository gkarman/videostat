package query

import "github.com/gkarman/demo/internal/application/blogger/query/view"

type VideoEnricher interface {
	Enrich(videos []*view.Video)
}
