package apigateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	_ "github.com/pitsanujiw/go-boilerplate/cmd/service/docs"
	"github.com/pitsanujiw/go-boilerplate/config"
	hhealth "github.com/pitsanujiw/go-boilerplate/internal/handler/health"
	hProxy "github.com/pitsanujiw/go-boilerplate/internal/handler/proxy"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
	"github.com/pitsanujiw/go-boilerplate/pkg/server"
)

// @title			Gateway API
// @version		1.0
// @description	Gateway provide information about activate of application.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func cliCommand(ctx context.Context) error {
	cfg, err := config.New("api_gateway_env.yaml")
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

	// initial handlers
	hhealth.New(srv, cfg)
	hProxy.New(srv, cfg)

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
		Name: "gateway",
		Aliases: []string{
			"gw",
		},
		Usage:       "running api-gateway",
		Description: "running api-gateway",
		Action: func(ctx *cli.Context) error {
			return cliCommand(ctx.Context)
		},
	}
}
