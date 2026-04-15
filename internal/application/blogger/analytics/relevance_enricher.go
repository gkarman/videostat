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
			// core retirement
			"retirement",
			"retire",
			"retiring",
			"early retirement",

			// pension
			"pension",
			"pension plan",
			"pension benefits",
			"defined benefit",
			"defined contribution",

			// social security (US/EU common topic)
			"social security",
			"social security benefits",
			"ss benefits",
			"retirement benefits",

			// taxes
			"retirement tax",
			"tax on retirement",
			"taxes in retirement",
			"pension tax",
			"tax deferred",
			"tax advantaged",

			// planning / finance context
			"retirement planning",
			"retirement savings",
			"401k",
			"ira",
			"roth ira",
			"pension fund",
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