package blogger

import "time"

type VideoGenerationStatus string

const (
	VideoGenerationPending    VideoGenerationStatus = "pending"
	VideoGenerationProcessing VideoGenerationStatus = "processing"
	VideoGenerationCompleted  VideoGenerationStatus = "completed"
	VideoGenerationFailed     VideoGenerationStatus = "failed"
)

type VideoGeneration struct {
	ID           string
	VideoID      string
	Platform     string
	ExternalID   string
	Status       VideoGenerationStatus
	S3URL        string
	ErrorMessage string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
