package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

func (p *RabbitPublisher) connect() error {
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

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

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
		conn.Close()
		return err
	}

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
