package events

import "time"

type VideoProcessingStartedV1 struct {
	EventType   string `json:"event_type"`
	EventID     string `json:"event_id"`
	VideoID     string `json:"video_id"`
	VideoURL    string `json:"video_url"`
	OccurredAt  time.Time
	PublishedAt time.Time
}
