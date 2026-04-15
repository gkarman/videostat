package analytics

import (
	"sort"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

type ViralEnricher struct {
	threshold float64
}

func NewViralEnricher() *ViralEnricher {
	return &ViralEnricher{
		threshold: 0.9, // top 10%
	}
}

func (e *ViralEnricher) Enrich(videos []*view.Video) {
	if len(videos) == 0 {
		return
	}

	grouped := e.groupByBlogger(videos)

	for _, group := range grouped {
		e.markViral(group)
	}
}

func (e *ViralEnricher) groupByBlogger(videos []*view.Video) map[string][]*view.Video {
	res := make(map[string][]*view.Video, len(videos))

	for _, v := range videos {
		res[v.BloggerURL] = append(res[v.BloggerURL], v)
	}

	return res
}

func (e *ViralEnricher) markViral(videos []*view.Video) {
	n := len(videos)
	if n == 0 {
		return
	}

	// engagement baseline
	first := e.engagement(videos[0])
	allSame := true

	for _, v := range videos[1:] {
		if e.engagement(v) != first {
			allSame = false
			break
		}
	}

	// если нет вариации — нет смысла выделять viral
	if allSame {
		return
	}

	sort.Slice(videos, func(i, j int) bool {
		return e.engagement(videos[i]) < e.engagement(videos[j])
	})

	cutoff := int(float64(n) * e.threshold)

	if cutoff >= n {
		return
	}

	for i := cutoff; i < n; i++ {
		videos[i].Viral = true
	}
}

func (e *ViralEnricher) engagement(v *view.Video) float64 {
	return float64(v.Views) +
		float64(v.Likes)*5 +
		float64(v.Comments)*10
}