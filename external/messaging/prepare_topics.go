package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func PrepareTopics(ctx context.Context, addr string, topics []string) error {
	conn, err := kafka.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if !isTopicAlreadyExists(partitions, topic) {
			if err := conn.CreateTopics(kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func isTopicAlreadyExists(partitions []kafka.Partition, topic string) bool {
	for _, p := range partitions {
		if p.Topic == topic {
			return true
		}
	}

	return false
}
