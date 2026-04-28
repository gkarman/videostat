package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapVideoSourceFound(e *blogger.VideoSourceFound) contracts.VideoSourceFoundV1 {
	return contracts.VideoSourceFoundV1{
		EventType:  contracts.EventVideoSourceFoundV1,
		EventID:    uuid.New().String(),
		VideoID:    e.VideoID,
		FileURL:    e.FileURL,
		OccurredAt: e.At,
	}
}
