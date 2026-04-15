package analytics

import (
	"sort"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

type VideoMetrics struct {
	viralThreshold float64
}

func NewVideoMetrics() *VideoMetrics {
	return &VideoMetrics{
		// top 10%
		viralThreshold: 0.9,
	}
}

func (s *VideoMetrics) Enrich(videos []*view.Video) {
	if len(videos) == 0 {
		return
	}

	byBlogger := s.groupByBlogger(videos)

	for _, group := range byBlogger {
		s.markViral(group)
	}
}

func (s *VideoMetrics) groupByBlogger(videos []*view.Video) map[string][]*view.Video {
	res := make(map[string][]*view.Video, len(videos))

	for _, v := range videos {
		res[v.BloggerURL] = append(res[v.BloggerURL], v)
	}

	return res
}

func (s *VideoMetrics) markViral(videos []*view.Video) {
	n := len(videos)
	if n == 0 {
		return
	}

	// 1. проверка на отсутствие разброса
	first := s.engagement(videos[0])
	allSame := true

	for _, v := range videos[1:] {
		if s.engagement(v) != first {
			allSame = false
			break
		}
	}

	if allSame {
		return // <- КЛЮЧЕВОЕ ИСПРАВЛЕНИЕ
	}

	// 2. сортировка
	sort.Slice(videos, func(i, j int) bool {
		return s.engagement(videos[i]) < s.engagement(videos[j])
	})

	// 3. percentile
	cutoff := int(float64(n) * s.viralThreshold)

	if cutoff >= n {
		return
	}

	for i := cutoff; i < n; i++ {
		videos[i].Viral = true
	}
}

func (s *VideoMetrics) engagement(v *view.Video) float64 {
	return float64(v.Views) +
		float64(v.Likes)*5 +
		float64(v.Comments)*10
}