package video

import "time"

type ContentStats struct {
	ContentID int64
	Views    *int64
	Likes    *int64
	Comments *int64
	Shares   *int64

	CollectedAt time.Time
}