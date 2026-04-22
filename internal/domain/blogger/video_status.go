package blogger

type VideoStatus string

const (
	VideoStatusCreated    VideoStatus = "created"
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
)

func (s VideoStatus) IsValid() bool {
	switch s {
	case
		VideoStatusCreated,
		VideoStatusProcessing,
		VideoStatusReady,
		VideoStatusFailed:
		return true
	default:
		return false
	}
}