package analytics

import (
	"testing"

	"github.com/gkarman/demo/internal/application/blogger/query/view"
)

func TestViralEnricher_Enrich(t *testing.T) {
	tests := []struct {
		name     string
		videos   []*view.Video
		assertFn func(t *testing.T, videos []*view.Video)
	}{
		{
			name: "outlier becomes viral inside single blogger",
			videos: []*view.Video{
				{BloggerURL: "b1", Views: 1000, Likes: 10, Comments: 5},
				{BloggerURL: "b1", Views: 1100, Likes: 12, Comments: 6},
				{BloggerURL: "b1", Views: 1200, Likes: 11, Comments: 7},
				{BloggerURL: "b1", Views: 1300, Likes: 10, Comments: 5},
				{BloggerURL: "b1", Views: 10000, Likes: 500, Comments: 200},
			},
			assertFn: func(t *testing.T, videos []*view.Video) {
				m := NewViralEnricher()
				m.Enrich(videos)

				var viral []*view.Video
				for _, v := range videos {
					if v.Viral {
						viral = append(viral, v)
					}
				}

				if len(viral) == 0 {
					t.Fatal("expected at least one viral video")
				}

				var best float64
				for _, v := range videos {
					e := float64(v.Views) + float64(v.Likes)*5 + float64(v.Comments)*10
					if e > best {
						best = e
					}
				}

				for _, v := range viral {
					e := float64(v.Views) + float64(v.Likes)*5 + float64(v.Comments)*10
					if e != best {
						t.Fatal("viral is not top engagement video")
					}
				}
			},
		},
		{
			name: "isolation between bloggers",
			videos: []*view.Video{
				{BloggerURL: "a", Views: 100, Likes: 5, Comments: 2},
				{BloggerURL: "a", Views: 110, Likes: 6, Comments: 3},
				{BloggerURL: "a", Views: 5000, Likes: 300, Comments: 100},

				{BloggerURL: "b", Views: 100000, Likes: 5000, Comments: 2000},
				{BloggerURL: "b", Views: 105000, Likes: 5200, Comments: 2100},
				{BloggerURL: "b", Views: 110000, Likes: 5300, Comments: 2200},
				{BloggerURL: "b", Views: 115000, Likes: 5400, Comments: 2300},
				{BloggerURL: "b", Views: 120000, Likes: 5500, Comments: 2400},
			},
			assertFn: func(t *testing.T, videos []*view.Video) {
				m := NewViralEnricher()
				m.Enrich(videos)

				var aViral, bViral bool

				for _, v := range videos {
					if v.Viral && v.BloggerURL == "a" {
						aViral = true
					}
					if v.Viral && v.BloggerURL == "b" {
						bViral = true
					}
				}

				if !aViral {
					t.Fatal("expected viral in blogger A")
				}

				if !bViral {
					t.Fatal("expected viral in blogger B")
				}
			},
		},
		{
			name: "no variance still produces no viral",
			videos: []*view.Video{
				{BloggerURL: "x", Views: 1000, Likes: 10, Comments: 5},
				{BloggerURL: "x", Views: 1000, Likes: 10, Comments: 5},
				{BloggerURL: "x", Views: 1000, Likes: 10, Comments: 5},
				{BloggerURL: "x", Views: 1000, Likes: 10, Comments: 5},
				{BloggerURL: "x", Views: 1000, Likes: 10, Comments: 5},
			},
			assertFn: func(t *testing.T, videos []*view.Video) {
				m := NewViralEnricher()
				m.Enrich(videos)

				for _, v := range videos {
					if v.Viral {
						t.Fatal("expected no viral in flat distribution")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertFn(t, tt.videos)
		})
	}
}