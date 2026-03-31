package mq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	cfg      Config
	queue    string
	bindings []string
	log      *slog.Logger
	backoff time.Duration
}

func NewRabbitConsumer(cfg Config, queue string, bindings []string, log *slog.Logger) *RabbitConsumer {
	return &RabbitConsumer{
		cfg:      cfg,
		queue:    queue,
		bindings: bindings,
		log:      log,
		backoff:  time.Second,
	}
}

func (c *RabbitConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		c.log.Info("rabbit consume loop start")

		err := c.consumeOnce(ctx, handler)
		if err != nil {
			wait := c.nextBackoff()

			c.log.Error("consume failed, will retry",
				"error", err,
				"retry_in", wait,
			)

			select {
			case <-time.After(wait):
				continue
			case <-ctx.Done():
				c.log.Info("consumer stopped by context")
				return ctx.Err()
			}
		}

		return nil
	}
}

func (c *RabbitConsumer) consumeOnce(ctx context.Context, handler func([]byte) error) error {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
	)

	c.log.Info("dial rabbit")

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	c.log.Info("rabbit connected")
	c.resetBackoff()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("qos: %w", err)
	}

	c.log.Info("declare queue", "queue", c.queue)

	_, err = ch.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	for _, key := range c.bindings {
		c.log.Info("bind queue", "key", key)

		if err := ch.QueueBind(
			c.queue,
			key,
			c.cfg.Exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("bind: %w", err)
		}
	}

	msgs, err := ch.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	c.log.Info("consumer started")

	closeCh := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("msgs channel closed")
			}

			c.log.Debug("message received", "routing_key", msg.RoutingKey)

			if err := handler(msg.Body); err != nil {
				c.log.Error("handler failed", "error", err)
				_ = msg.Nack(false, true)
				continue
			}

			_ = msg.Ack(false)

		case err := <-closeCh:
			return fmt.Errorf("channel closed: %v", err)

		case <-ctx.Done():
			c.log.Info("consumeOnce stopped by context")
			return ctx.Err()
		}
	}
}

func (c *RabbitConsumer) nextBackoff() time.Duration {
	max := 30 * time.Second

	d := c.backoff
	c.backoff *= 2

	if c.backoff > max {
		c.backoff = max
	}

	return d
}

func (c *RabbitConsumer) resetBackoff() {
	c.backoff = time.Second
}