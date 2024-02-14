package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"

	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

type KafkaConsumer interface {
	Subscriber(ctx context.Context, fn ConsumerReceiveFunc)
	Close() error
}

type kafkaConsumer struct {
	consumer *kafka.Reader
	logger   *log.Logger
	topic    string
}

type ConsumerReceiveFunc func(ctx context.Context, msg kafka.Message) error

const (
	defaultRetry = 10              // retry 10 rounds
	defaultDelay = 5 * time.Second // 5 sec and exponential time * 10 time
)

func NewConsumer(logger *log.Logger, topic, groupID string, addr ...string) KafkaConsumer {
	c := kafka.NewReader(kafka.ReaderConfig{
		MaxAttempts:           25,
		Brokers:               addr,
		Topic:                 topic,
		GroupID:               groupID,
		MaxBytes:              10e6, // 10MB
		MaxWait:               100 * time.Millisecond,
		OffsetOutOfRangeError: true,
	})

	return &kafkaConsumer{
		consumer: c,
		logger:   logger,
		topic:    topic,
	}
}

func (k *kafkaConsumer) Close() error {
	return k.consumer.Close()
}

func (k *kafkaConsumer) Subscriber(ctx context.Context, fn ConsumerReceiveFunc) {
	for {
		msg, err := k.consumer.FetchMessage(ctx)
		if err != nil {
			k.logger.Sugar().With("topic", k.topic).Errorln(fmt.Sprintf("unable to read message: %s", err))
			break
		}

		k.logger.Sugar().With(
			"topic", k.topic,
			"value", string(msg.Value),
			"key", string(msg.Key),
			"partition", msg.Partition,
			"time", msg.Time.String(),
			"offset", msg.Offset,
		).Infoln(fmt.Sprintf("[KAFKA-RECEIVED] topic: %s time: %s", k.topic, msg.Time.String()))

		go func() {
			if err := retry.Do(
				func() error {
					if err := fn(ctx, msg); err != nil {
						k.logger.Sugar().With("topicId", k.topic).Errorln(fmt.Sprintf("[KAFKA-RECEIVE-FAILED] failed message: %s", err))

						return err
					}
					if err := k.consumer.CommitMessages(ctx, msg); err != nil {
						k.logger.Sugar().With("topicId", k.topic).Errorln(fmt.Sprintf("[KAFKA-RECEIVE-FAILED] failed to commit message: %s", err))

						return err
					}

					return nil
				},
				retry.Attempts(defaultRetry),
				retry.Delay(defaultDelay),
				retry.DelayType(retry.BackOffDelay),
			); err != nil {
				k.logger.Sugar().With("topicId", k.topic).Errorln(fmt.Sprintf("[KAFKA-RECEIVE-FAILED] failed to retries to run subscribe: %s", err))

				if err := k.consumer.CommitMessages(ctx, msg); err != nil {
					k.logger.Sugar().With("topicId", k.topic).Errorln(fmt.Sprintf("[KAFKA-RECEIVE-FAILED] failed to commit message after failed to retries: %s", err))
				}
			}

			k.logger.Sugar().With(
				"topic", k.topic,
				"value", string(msg.Value),
				"key", string(msg.Key),
				"partition", msg.Partition,
				"time", msg.Time.String(),
				"offset", msg.Offset,
			).Infoln(fmt.Sprintf("[KAFKA-RECEIVED] topic: %s time: %s is successful", k.topic, msg.Time.String()))
		}()
	}
}
