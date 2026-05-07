package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoPromptGenerated(e *blogger.VideoPromptGenerated) contracts.VideoPromptGeneratedV1 {
	return contracts.VideoPromptGeneratedV1{
		EventType:  contracts.EventVideoPromptGeneratedV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		OccurredAt: e.At,
	}
}
