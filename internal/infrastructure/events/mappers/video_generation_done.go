package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoGenerationDone(e *blogger.VideoGenerationDone) contracts.VideoGenerationDoneV1 {
	return contracts.VideoGenerationDoneV1{
		EventType:  contracts.EventVideoGenerationDoneV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		ChatID:     e.ChatID,
		S3URL:      e.S3URL,
		OccurredAt: e.At,
	}
}

func MapVideoGenerationError(e *blogger.VideoGenerationError) contracts.VideoGenerationErrorV1 {
	return contracts.VideoGenerationErrorV1{
		EventType:  contracts.EventVideoGenerationErrorV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		ChatID:     e.ChatID,
		Reason:     e.Reason,
		OccurredAt: e.At,
	}
}
