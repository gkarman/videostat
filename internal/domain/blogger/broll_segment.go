package blogger

import "time"

const (
	BrollStatusPending    = "pending"
	BrollStatusProcessing = "processing"
	BrollStatusDone       = "done"
	BrollStatusFailed     = "failed"
)

type BrollSegment struct {
	ID                   string
	VideoID              string
	Position             int
	StartMS              int
	EndMS                int
	Text                 string
	BrollPrompt          string
	BrollURL             *string
	GenerationStatus     string
	GenerationExternalID *string
	GenerationError      *string
	CreatedAt            time.Time
}
