package events

import "time"

type VideoCompositionDoneV1 struct {
	EventType  string    `json:"event_type"`
	EventID    string    `json:"event_id"`
	VideoID    string    `json:"video_id"`
	ChatID     int64     `json:"chat_id"`
	ResultURL  string    `json:"result_url"`
	SourceURL  string    `json:"source_url"`
	OccurredAt time.Time `json:"occurred_at"`
}
