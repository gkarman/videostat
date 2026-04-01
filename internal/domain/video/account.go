package video

import "time"



type Account struct {
	ID         int64
	PlatformID int16
	ExternalID string
	Title      string
	URL        string
	CreatedAt time.Time
	UpdatedAt time.Time
}
