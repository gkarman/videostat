package video

import "time"

type Content struct {
	ID          int64
	AccountID   int64
	ExternalID  string
	Title       string
	PublishedAt *time.Time
	DurationSec *int
	CreatedAt   time.Time
}
