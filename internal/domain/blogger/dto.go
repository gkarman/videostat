package blogger

import "time"

type CreateBloggerDto struct {
	ID       string
	PlatformID int
	URL      string
}

type CreateVideoDto struct {
	ID          string
	BloggerID   string
	ExternalID  string
	URL         string
	Title       string

	Views    int64
	Likes    int64
	Comments int64

	PublishedAt time.Time
}
