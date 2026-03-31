package events

import "time"

type CarUpdatedV1 struct {
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
	CarID string `json:"car_id"`
	NameOld  string `json:"name_old"`
	NameNew  string `json:"name_new"`
	OccurredAt  time.Time
	PublishedAt time.Time
}
