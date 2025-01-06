//go:generate mockery --name=PublisherInterface --output=../util/mocks --outpkg=mocks
package event

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type SubscriberInterface interface {
	Consume(ctx context.Context, handler func(key, value []byte) error)
	Close() error
}

type Subscriber struct {
	reader *kafka.Reader
}

func NewSubscriber(broker string, topic string, groupID string) *Subscriber {
	return &Subscriber{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Subscriber) Consume(ctx context.Context, handler func(key, value []byte) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}

		if err := handler(msg.Key, msg.Value); err != nil {
			log.Printf("failed to process message: %v", err)
		}
	}
}

func (c *Subscriber) Close() error {
	return c.reader.Close()
}
