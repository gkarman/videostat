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

	events []any
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

func (v *Video) MarkStartProcessing() error {
	err := v.ChangeStatus(VideoStatusProcessing)
	if err != nil {
		return err
	}

	v.ErrorStage = nil
	v.ErrorMessage = nil

	v.addEvent(&VideoProcessingStarted{
		VideoID:  v.ID,
		VideoURL: v.URL,
		At:       time.Now(),
	})

	return nil
}

func (v *Video) MarkFailProcessing(stage VideoErrorStage, err error) error {
	if err := v.ChangeStatus(VideoStatusFailed); err != nil {
		return err
	}
	v.ErrorStage = &stage
	msg := err.Error()
	v.ErrorMessage = &msg

	v.addEvent(&VideoProcessingFailed{
		VideoID:      v.ID,
		Stage:        stage,
		ErrorMessage: msg,
		At:           time.Now(),
	})

	return nil
}

func (v *Video) SourceFound(fileURL string) {
	v.addEvent(&VideoSourceFound{
		VideoID: v.ID,
		FileURL: fileURL,
		At:      time.Now(),
	})
}

func (v *Video) AnalyzeDone() {
	v.addEvent(&VideoAnalyzeDone{
		VideoID: v.ID,
		At:      time.Now(),
	})
}

func (v *Video) PromptGenerated() {
	v.addEvent(&VideoPromptGenerated{
		VideoID: v.ID,
		At:      time.Now(),
	})
}

func (v *Video) MarkReady() error {
	return v.ChangeStatus(VideoStatusReady)
}

func (v *Video) ChangeStatus(to VideoStatus) error {
	if !v.canChangeStatusTo(to) {
		return errors.New("invalid status transition")
	}

	v.Status = to
	return nil
}

func (v *Video) canChangeStatusTo(to VideoStatus) bool {
	switch v.Status {

	case VideoStatusCreated:
		return to == VideoStatusProcessing || to == VideoStatusFailed

	case VideoStatusProcessing:
		return to == VideoStatusReady || to == VideoStatusFailed

	case VideoStatusReady:
		return false

	case VideoStatusFailed:
		return to == VideoStatusProcessing

	default:
		return false
	}
}

func (v *Video) addEvent(e any) {
	v.events = append(v.events, e)
}

func (v *Video) PullEvents() []any {
	evs := v.events
	v.events = nil
	return evs
}
