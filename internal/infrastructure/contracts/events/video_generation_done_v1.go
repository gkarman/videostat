package events

import "time"

type VideoGenerationDoneV1 struct {
	EventType  string    `json:"event_type"`
	EventID    string    `json:"event_id"`
	VideoID    string    `json:"video_id"`
	ChatID     int64     `json:"chat_id"`
	S3URL      string    `json:"s3_url"`
	OccurredAt time.Time `json:"occurred_at"`
}
