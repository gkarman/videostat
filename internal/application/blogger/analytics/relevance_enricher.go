package analytics

import (
	"strings"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

type RelevanceEnricher struct {
	keywords []string
}

func NewRelevanceEnricher() *RelevanceEnricher {
	return &RelevanceEnricher{
		keywords: []string{
			"пенсия",
			"пенсион",
			"пенсионный возраст",
			"pension",
			"retirement",
		},
	}
}

func (e *RelevanceEnricher) Enrich(videos []*view.Video) {
	for _, v := range videos {
		v.IsRelevant = e.isRelevant(v.Title)
	}
}

func (e *RelevanceEnricher) isRelevant(title string) bool {
	t := strings.ToLower(title)

	for _, k := range e.keywords {
		if strings.Contains(t, k) {
			return true
		}
	}

	return false
}