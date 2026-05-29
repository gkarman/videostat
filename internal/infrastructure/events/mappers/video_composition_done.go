package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoCompositionDone(e *blogger.VideoCompositionDone) contracts.VideoCompositionDoneV1 {
	return contracts.VideoCompositionDoneV1{
		EventType:  contracts.EventVideoCompositionDoneV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		ChatID:     e.ChatID,
		ResultURL:  e.ResultURL,
		SourceURL:  e.SourceURL,
		OccurredAt: e.At,
	}
}
