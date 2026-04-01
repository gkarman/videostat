package video

import "time"

type AccountStats struct {
	AccountID int64
	Followers  *int64
	TotalViews *int64

	CollectedAt time.Time
}