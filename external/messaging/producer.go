package messaging

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

type KafkaProducer interface {
	PublishMessages(msg kafka.Message) error
	Close() error
}

type kafkaProducer struct {
	producer *kafka.Writer
	logger   *log.Logger
	topic    string
}

func NewProducer(logger *log.Logger, topic string, addr ...string) KafkaProducer {
	p := &kafka.Writer{
		Addr:         kafka.TCP(addr...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}

	return &kafkaProducer{
		producer: p,
		topic:    topic,
		logger:   logger,
	}
}

func (k *kafkaProducer) Close() error {
	return k.producer.Close()
}

func (k *kafkaProducer) PublishMessages(msg kafka.Message) error {
	ctx := context.Background()

	if err := k.producer.WriteMessages(ctx, msg); err != nil {
		k.logger.Sugar().With("topic", k.topic).Errorln(fmt.Sprintf("[PUBLISH ERROR] failed to publish %v to topic %v : %v", msg.Partition, k.topic, err))

		return err
	}

	k.logger.Sugar().With(
		"topic", k.topic,
		"value", string(msg.Value),
		"key", string(msg.Key),
		"partition", msg.Partition,
		"time", msg.Time.String(),
		"offset", msg.Offset,
	).Infoln(fmt.Sprintf("[PUBLISHED] partition %v to topic %v", msg.Partition, k.topic))

	return nil
}
