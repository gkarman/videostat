package events

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	car_handlers "github.com/gkarman/demo/internal/application/car/handlers"
	blogger_handlers "github.com/gkarman/demo/internal/application/blogger/handlers"
	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/domain/car"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
)

func RegisterEventHandlers(d *dispatcher.Dispatcher, log *slog.Logger, publisher application.Publisher) {
	d.Register(&car.Created{}, car_handlers.CarCreatedLogHandler(log))
	d.Register(&car.Created{}, car_handlers.CarCreatedToRabbitHandler(publisher, log))
	d.Register(&blogger.Created{}, blogger_handlers.BloggerCreatedToRabbitHandler(publisher, log))
}
