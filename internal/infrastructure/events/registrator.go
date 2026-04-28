package events

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	blogger_handlers "github.com/gkarman/demo/internal/application/blogger/handlers"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
)

func RegisterEventHandlers(d *dispatcher.Dispatcher, log *slog.Logger, publisher application.Publisher) {
	d.Register(&blogger.Created{}, blogger_handlers.BloggerCreatedToRabbitHandler(publisher, log))
	d.Register(&blogger.VideoProcessingStarted{}, blogger_handlers.VideoProcessingStartedToRabbitHandler(publisher, log))
	d.Register(&blogger.VideoSourceFound{}, blogger_handlers.VideoSourceFoundToRabbitHandler(publisher, log))
	d.Register(&blogger.VideoAnalyzeDone{}, blogger_handlers.VideoAnalyzeDoneToRabbitHandler(publisher, log))
}
