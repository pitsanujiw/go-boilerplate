package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/urfave/cli/v2"

	_ "github.com/pitsanujiw/go-boilerplate/cmd/service/docs"
	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/external/messaging"
	hhealth "github.com/pitsanujiw/go-boilerplate/internal/handler/health"
	"github.com/pitsanujiw/go-boilerplate/pkg/database"
	"github.com/pitsanujiw/go-boilerplate/pkg/database/gen"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
	"github.com/pitsanujiw/go-boilerplate/pkg/server"
)

// @title			Service API
// @version		1.0
// @description	Service provide information about activate of application.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func cliCommand(ctx context.Context) error {
	cfg, err := config.New("env.yaml")
	if err != nil {
		return fmt.Errorf("unable to load app config: %v", err)
	}

	logger, err := log.New(cfg)
	if err != nil {
		return fmt.Errorf("unable to initial logger: %v", err)
	}

	// Server
	srv, err := server.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("unable to initial server: %v", err)
	}

	if err := messaging.PrepareTopics(ctx, cfg.Kafka.Addr[0], []string{cfg.Kafka.Producer.TopicID}); err != nil {
		logger.Sugar().
			With("service", cfg.Name, "time", time.Now().Format(time.RFC3339), "tz", cfg.TZ).
			Errorln(fmt.Sprintf("cannot prepare topic err: %v", err))

		return err
	}

	// example for kafka producer external
	p := messaging.NewProducer(logger, cfg.Kafka.Producer.TopicID, cfg.Kafka.Addr...)
	defer p.Close()

	go func() {
		for i := 0; i < 5; i++ {
			if err := p.PublishMessages(kafka.Message{
				Value: []byte(fmt.Sprintf("Topic-1: %v", i)),
			}); err != nil {
				continue
			}
		}
	}()

	pool, err := database.AcquireDBPool(ctx, cfg.Database)
	if err != nil {
		return err
	}

	defer pool.Close()

	g := gen.New()
	data, err := g.Get(ctx, pool)
	if err != nil {
		return err
	}

	logger.Sugar().Infoln(data)

	hhealth.New(srv, cfg)

	go func() {
		logger.Sugar().
			With(
				"service", cfg.Name,
				"tz", cfg.TZ,
			).
			Info(fmt.Sprintf("starting member activate time: %s",
				time.Now().Format(time.RFC3339)))

		switch err := srv.Start(ctx); {
		case errors.Is(err, http.ErrServerClosed):
			logger.Error("Server shutting down...")
		case err != nil:
			logger.Error("Server shutting down unexpectedly...")
		}
	}()

	// below code allows for graceful shut down
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := srv.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	logger.Sugar().
		With("service", cfg.Name, "time", time.Now().Format(time.RFC3339), "tz", cfg.TZ).
		Warnln("Shutdown Complete")

	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name: "service",
		Aliases: []string{
			"srv",
		},
		Usage:       "running service",
		Description: "running service",
		Action: func(ctx *cli.Context) error {
			return cliCommand(ctx.Context)
		},
	}
}
