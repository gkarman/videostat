package blogger

import "time"

type VideoWatcher struct {
	ID        string
	VideoID   string
	ChatID    int64
	CreatedAt time.Time
}
