package notify

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gkarman/demo/internal/infrastructure/mq"
)


type Worker struct {
	log      *slog.Logger
	consumer *mq.RabbitConsumer
	router   *Router
}

func New(log *slog.Logger, consumer *mq.RabbitConsumer, router *Router) *Worker {
	return &Worker{
		log:      log,
		consumer: consumer,
		router:   router,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("worker_notify started")

	go func() {
		err := w.consumer.Consume(ctx, w.handleMessage)
		if err != nil {
			w.log.Error("consumer stopped", "error", err)
		}
	}()

	<-ctx.Done()
	w.log.Info("worker_notify shutting down", "reason", ctx.Err())

	return nil
}

func (w *Worker) handleMessage(body []byte) error {
	var base struct {
		EventType string `json:"event_type"`
	}

	if err := json.Unmarshal(body, &base); err != nil {
		return err
	}

	return w.router.Handle(base.EventType, body)
}
