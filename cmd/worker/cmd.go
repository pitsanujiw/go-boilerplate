package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/urfave/cli/v2"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/external/messaging"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

func cliCommand(ctx context.Context) error {
	cfg, err := config.New("worker_env.yaml")
	if err != nil {
		return err
	}

	logger, err := log.New(cfg)
	if err != nil {
		return err
	}

	if err := messaging.PrepareTopics(ctx, cfg.Kafka.Addr[0], cfg.Kafka.Consumer.TopicIDs); err != nil {
		logger.Sugar().
			With("service", cfg.Name, "time", time.Now().Format(time.RFC3339), "tz", cfg.TZ).
			Errorln(fmt.Sprintf("cannot prepare topic err: %v", err))

		return err
	}

	gws := messaging.NewConsumerGateway()
	// Registers consumer to support running consumer by pods
	gws.Registers(
		messaging.ConsumerGatewayRegister{
			TopicID: messaging.ExampleTopicID,
			Consumer: func(ctx context.Context, msg kafka.Message) error {
				logger.Sugar().
					With("topic", msg.Topic).Infoln(string(msg.Value))

				return nil
			},
		},
	)

	for _, topicID := range cfg.Kafka.Consumer.TopicIDs {
		consumer, err := gws.Get(messaging.ConsumerTopicID(topicID))
		if err != nil {
			return err
		}

		// example for consumer to read data from kafka
		c := messaging.NewConsumer(logger, topicID, cfg.Kafka.Consumer.GroupID, cfg.Kafka.Addr...)
		defer c.Close()
		// spawn a other pod to subscriber my consumer messages
		go c.Subscriber(ctx, consumer)
	}

	logger.Sugar().With(
		"service", cfg.Name,
		"time", time.Now().String(),
	).Infoln("startup worker")

	// below code allows for graceful shut down
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown

	logger.Sugar().With(
		"service", cfg.Name,
		"time", time.Now().String(),
	).Warnln("Shutdown Complete")

	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name: "worker",
		Aliases: []string{
			"w",
		},
		Usage:       "running member activate worker service",
		Description: "running member activate worker service",
		Action: func(ctx *cli.Context) error {
			return cliCommand(ctx.Context)
		},
	}
}
