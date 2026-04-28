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

type VideoSourceFound struct {
	VideoID string
	FileURL string
	At      time.Time
}

type VideoAnalyzeDone struct {
	VideoID string
	At      time.Time
}

type VideoProcessingFailed struct {
	VideoID      string
	Stage        VideoErrorStage
	ErrorMessage string
	At           time.Time
}
