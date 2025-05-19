package xmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer[T any] interface {
	Produce(ctx context.Context, msg T) error
	Close()
}

var _ Producer[any] = (*GeneralProducer[any])(nil)

type GeneralProducer[T any] struct {
	producer *kafka.Producer
	topic    string
}

func (g *GeneralProducer[T]) Produce(ctx context.Context, msg T) error {
	data, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	// send messages
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := g.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &g.topic,
				Partition: kafka.PartitionAny,
			},
			Value: data,
		}, deliveryChan)

		if err != nil {
			var kafkaErr kafka.Error
			if ok := errors.As(err, &kafkaErr); ok && kafkaErr.Code() == kafka.ErrQueueFull {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second):
					continue
				}
			}
			return fmt.Errorf("send message to topic [%s] failed: %w", g.topic, err)
		}

		for {
			flushed := g.producer.Flush(int(time.Second))
			if flushed == 0 {
				break
			}
		}
		break
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-deliveryChan:
		switch event := e.(type) {
		case *kafka.Message:
			if event.TopicPartition.Error != nil {
				return fmt.Errorf("failed to deliver message: %w", event.TopicPartition.Error)
			}
			return nil
		case *kafka.Error:
			return fmt.Errorf("kafka error: %w", event)
		default:
			return fmt.Errorf("unknown kafka event type: %w", errors.New(event.String()))
		}
	}
}

func (g *GeneralProducer[T]) Close() {
	for {
		flushed := g.producer.Flush(int(time.Second))
		if flushed == 0 {
			break
		}
	}
	g.producer.Close()
}

func NewGeneralProducer[T any](topic string, producer *kafka.Producer) *GeneralProducer[T] {
	return &GeneralProducer[T]{
		producer: producer,
		topic:    topic,
	}
}
