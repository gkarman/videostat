package events

import "time"

type BloggerCreatedV1 struct {
	EventType   string `json:"event_type"`
	EventID     string `json:"event_id"`
	BloggerID   string `json:"blogger_id"`
	BloggerURL  string `json:"blogger_url"`
	OccurredAt  time.Time
	PublishedAt time.Time
}
