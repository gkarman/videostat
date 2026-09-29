package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RabbitPublisher struct {
	cfg    Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	mu     sync.RWMutex
	closed bool
}

func NewRabbitPublisher(cfg Config) (*RabbitPublisher, error) {
	p := &RabbitPublisher{cfg: cfg}

	if err := p.connect(); err != nil {
		return nil, err
	}

	go p.reconnectLoop()

	return p, nil
}

func (p *RabbitPublisher) connect() (err error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		p.cfg.User,
		p.cfg.Password,
		p.cfg.Host,
		p.cfg.Port,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = ch.Close()
		}
	}()

	err = ch.ExchangeDeclare(
		p.cfg.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// С этого момента ресурсы "переданы" в publisher
	p.mu.Lock()
	p.conn = conn
	p.ch = ch
	p.mu.Unlock()

	return nil
}

func (p *RabbitPublisher) reconnectLoop() {
	for {
		p.mu.RLock()
		conn := p.conn
		p.mu.RUnlock()

		if conn == nil {
			time.Sleep(p.cfg.ReconnectDelay)
			continue
		}

		closeCh := make(chan *amqp.Error)
		conn.NotifyClose(closeCh)

		err := <-closeCh
		if err != nil {
			fmt.Println("rabbit disconnected:", err)
		}

		for {
			if p.isClosed() {
				return
			}

			fmt.Println("reconnecting...")

			if err := p.connect(); err != nil {
				time.Sleep(p.cfg.ReconnectDelay)
				continue
			}

			fmt.Println("reconnected")
			break
		}
	}
}

func (p *RabbitPublisher) Publish(ctx context.Context, key string, body []byte) error {
	// Спан публикации — дочерний к текущей операции (команда в боте, HTTP-запрос, обработка события).
	ctx, span := tracer.Start(ctx, "publish "+key,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "rabbitmq"),
			attribute.String("messaging.destination.name", p.cfg.Exchange),
			attribute.String("messaging.rabbitmq.destination.routing_key", key),
		),
	)
	defer span.End()

	// traceparent (id трейса + id этого спана) — в заголовки сообщения.
	// Consumer прочитает его и продолжит тот же трейс.
	headers := amqp.Table{}
	otel.GetTextMapPropagator().Inject(ctx, headersCarrier(headers))

	err := p.publish(ctx, key, body, headers)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (p *RabbitPublisher) publish(ctx context.Context, key string, body []byte, headers amqp.Table) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher closed")
	}

	if p.ch == nil {
		return fmt.Errorf("no channel")
	}

	return p.ch.PublishWithContext(
		ctx,
		p.cfg.Exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Headers:     headers,
			Body:        body,
		},
	)
}

func (p *RabbitPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true

	if p.ch != nil {
		_ = p.ch.Close()
	}

	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}

func (p *RabbitPublisher) isClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.closed
}
