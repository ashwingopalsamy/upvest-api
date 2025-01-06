//go:generate mockery --name=PublisherInterface --output=../util/mocks --outpkg=mocks
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type PublisherInterface interface {
	Publish(ctx context.Context, key, value []byte) error
	Close() error
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(broker string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
