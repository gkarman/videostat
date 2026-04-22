package blogger

import (
	"errors"
	"time"
)

type Video struct {
	ID         string
	BloggerID  string
	ExternalID string
	URL        string
	Title      string
	Views      int64
	Likes      int64
	Comments   int64

	Status       VideoStatus
	ErrorStage   *VideoErrorStage
	ErrorMessage *string

	PublishedAt time.Time
	CreatedAt   time.Time
}

func NewVideo(dto CreateVideoDto) *Video {
	return &Video{
		ID:          dto.ID,
		BloggerID:   dto.BloggerID,
		ExternalID:  dto.ExternalID,
		URL:         dto.URL,
		Title:       dto.Title,
		Views:       dto.Views,
		Likes:       dto.Likes,
		Comments:    dto.Comments,
		PublishedAt: dto.PublishedAt,
		CreatedAt:   time.Now(),
		Status:      VideoStatusCreated,
	}
}


func (v *Video) ChangeStatus(to VideoStatus) error {
	if v.Status == to {
		return nil
	}

	if !v.СanChangeStatusTo(to) {
		return errors.New("invalid status transition")
	}

	v.Status = to
	return nil
}

func (v *Video) СanChangeStatusTo(to VideoStatus) bool {
	switch v.Status {

	case VideoStatusCreated:
		return to == VideoStatusProcessing || to == VideoStatusFailed

	case VideoStatusProcessing:
		return to == VideoStatusReady || to == VideoStatusFailed

	case VideoStatusReady:
		return false

	case VideoStatusFailed:
		return false

	default:
		return false
	}
}
