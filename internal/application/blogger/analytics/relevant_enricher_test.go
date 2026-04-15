package analytics

import (
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

func TestRelevanceEnricher(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected bool
	}{
		{
			name:     "retirement planning with 401k and ira",
			title:    "How to plan retirement with 401k and IRA",
			expected: true,
		},
		{
			name:     "retirement general topic",
			title:    "Retirement planning basics",
			expected: true,
		},
		{
			name:     "pension benefits and taxes",
			title:    "Understanding pension benefits and tax rules",
			expected: true,
		},
		{
			name:     "social security explanation",
			title:    "Social security benefits explained",
			expected: true,
		},
		{
			name:     "non relevant cooking content",
			title:    "How to cook pasta like a chef",
			expected: false,
		},
		{
			name:     "marketing unrelated content",
			title:    "Best marketing strategies for 2026 growth",
			expected: false,
		},
	}

	enricher := NewRelevanceEnricher()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &view.Video{
				Title: tt.title,
			}

			enricher.Enrich([]*view.Video{v})

			if v.IsRelevant != tt.expected {
				t.Fatalf(
					"expected %v, got %v (title=%s)",
					tt.expected,
					v.IsRelevant,
					tt.title,
				)
			}
		})
	}
}