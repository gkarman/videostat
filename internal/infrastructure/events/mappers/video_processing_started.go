package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoProcessingStarted(e *blogger.VideoProcessingStarted) contracts.VideoProcessingStartedV1 {
	return contracts.VideoProcessingStartedV1{
		EventType:  contracts.EventVideoProcessingStartedV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		VideoURL:   e.VideoURL,
		OccurredAt: e.At,
	}
}
