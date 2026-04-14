package mappers

import (
	"github.com/gkarman/demo/internal/domain/blogger"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/google/uuid"
)

func MapBloggerCreated(e *blogger.Created) contracts.BloggerCreatedV1 {
	return contracts.BloggerCreatedV1{
		EventType:  contracts.EventCarCreatedV1,
		EventID:    uuid.New().String(),
		BloggerID:  e.ID,
		BloggerURL: e.URL,
		OccurredAt: e.At,
	}
}
