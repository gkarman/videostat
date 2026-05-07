package events

import "time"

type VideoPromptGeneratedV1 struct {
	EventType   string    `json:"event_type"`
	EventID     string    `json:"event_id"`
	VideoID     string    `json:"video_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	PublishedAt time.Time `json:"published_at"`
}
