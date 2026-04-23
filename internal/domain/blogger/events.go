package blogger

import "time"

type Created struct {
	ID  string
	URL string
	At  time.Time
}

type VideoProcessingStarted struct {
	VideoID  string
	At  time.Time
}
