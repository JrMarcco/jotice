package xmq

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer interface {
	Read(timeout time.Duration) (*kafka.Message, error)
	Commit(m *kafka.Message) ([]kafka.TopicPartition, error)

	Pause(partitions []kafka.TopicPartition) error
	Resume(partitions []kafka.TopicPartition) error

	Poll(timeoutMs int32) error
	Assignment() ([]kafka.TopicPartition, error)
}
