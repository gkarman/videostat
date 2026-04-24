package blogger

import "time"

type Created struct {
	ID  string
	URL string
	At  time.Time
}

type VideoProcessingStarted struct {
	VideoID  string
	VideoURL string
	At       time.Time
}
