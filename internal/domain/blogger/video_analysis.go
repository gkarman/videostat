package blogger

import "time"

type VideoAnalysis struct {
	ID         string
	VideoID    string
	Provider   string
	RawPayload []byte
	CreatedAt  time.Time
}
