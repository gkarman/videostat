package mappers

import (
	"github.com/gkarman/demo/internal/domain/car"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapCarCreated(e *car.Created) contracts.CarCreatedV1 {
	return contracts.CarCreatedV1{
		EventType:  contracts.EventCarCreatedV1,
		EventID:    uuid.New().String(),
		CarID:      e.ID,
		Name:       e.Name,
		OccurredAt: e.At,
	}
}
