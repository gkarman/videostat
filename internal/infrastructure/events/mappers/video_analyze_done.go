package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoAnalyzeDone(e *blogger.VideoAnalyzeDone) contracts.VideoAnalyzeDoneV1 {
	return contracts.VideoAnalyzeDoneV1{
		EventType:  contracts.EventVideoAnalyzeDoneV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		OccurredAt: e.At,
	}
}
