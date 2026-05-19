package blogger

type VideoStatus string

const (
	VideoStatusCreated              VideoStatus = "created"
	VideoStatusProcessing           VideoStatus = "processing"
	VideoStatusGenerationProcessing VideoStatus = "generation_processing"
	VideoStatusGenerationFailed     VideoStatus = "generation_failed"
	VideoStatusReady                VideoStatus = "ready"
	VideoStatusFailed               VideoStatus = "failed"
)

func (s VideoStatus) IsValid() bool {
	switch s {
	case
		VideoStatusCreated,
		VideoStatusProcessing,
		VideoStatusGenerationProcessing,
		VideoStatusGenerationFailed,
		VideoStatusReady,
		VideoStatusFailed:
		return true
	default:
		return false
	}
}