package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

type httpServer struct {
	app       *fiber.App
	cfg       *config.App
	log       *log.Logger
	api       fiber.Router
	publicApi fiber.Router
}

type HTTPServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Server() *fiber.App
	PrivateEndpoint() fiber.Router
	PublicEndpoint() fiber.Router
}

func New(cfg *config.App, log *log.Logger) (HTTPServer, error) {
	var (
		app = fiber.New(fiber.Config{
			Prefork: false,
			AppName: cfg.Name,
		})
		// Create new group for API
		basePath = cfg.BasePath
		apiPath  = "/api/v1"
	)

	app.Use(
		logger.New(logger.Config{
			Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
			TimeFormat: "2 Jan 2006 15:04:05",
			TimeZone:   cfg.TZ,
		}),
		cors.New(),
		idempotency.New(),
		requestid.New(),
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		helmet.New(),
	)
	var (
		public = app.Group(basePath)
		api    = app.Group(fmt.Sprint(basePath, apiPath))
	)

	app.Get("/sw/*", swagger.HandlerDefault) // default
	// Initialize default config (Assign the middleware to /metrics)
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: cfg.Name,
	}))

	return &httpServer{
		app:       app,
		api:       api,
		publicApi: public,
		log:       log,
		cfg:       cfg,
	}, nil
}

func (svr *httpServer) Start(ctx context.Context) error {
	svr.log.Sugar().With("service", svr.cfg.Name, "tz", svr.cfg.TZ).Infoln("Start server")

	if err := svr.app.Listen(fmt.Sprintf(":%s", svr.cfg.Port)); err != nil {
		return err
	}

	return nil
}

func (svr *httpServer) Stop(ctx context.Context) error {
	svr.log.Sugar().With("service", svr.cfg.Name, "tz", svr.cfg.TZ).Infoln("Shutdown server")

	return svr.app.Shutdown()
}

func (svr *httpServer) Server() *fiber.App {
	return svr.app
}

func (svr *httpServer) PrivateEndpoint() fiber.Router {
	return svr.api
}

func (svr *httpServer) PublicEndpoint() fiber.Router {
	return svr.publicApi
}
