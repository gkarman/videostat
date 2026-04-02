package blogger

import "time"

type Video struct {
	ID      int
	Title   string
	URL     string
	Views   int
	Likes   int
	PublishedAt time.Time
}
