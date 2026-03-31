package events

import "time"

type CarCreatedV1 struct {
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
	CarID string `json:"car_id"`
	Name  string `json:"name"`
	OccurredAt  time.Time
	PublishedAt time.Time
}
