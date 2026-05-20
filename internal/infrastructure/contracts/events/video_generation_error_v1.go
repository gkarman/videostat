package events

import "time"

type VideoGenerationErrorV1 struct {
	EventType  string    `json:"event_type"`
	EventID    string    `json:"event_id"`
	VideoID    string    `json:"video_id"`
	ChatID     int64     `json:"chat_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}
