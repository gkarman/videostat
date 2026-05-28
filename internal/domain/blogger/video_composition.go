package blogger

import "time"

const (
	CompositionStatusProcessing = "processing"
	CompositionStatusDone       = "done"
	CompositionStatusFailed     = "failed"
)

type VideoComposition struct {
	ID         string
	VideoID    string
	Status     string
	ExternalID *string
	ResultURL  *string
	Error      *string
	CreatedAt  time.Time
}
