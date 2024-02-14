package server

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

type TestServer struct {
	cfg       *config.App
	logger    *log.Logger
	api       fiber.Router
	publicApi fiber.Router
	app       *fiber.App
}

func NewTest(logger *log.Logger) (HTTPServer, error) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	api := app.Group("/api/v1")
	return &TestServer{
		cfg:       nil,
		logger:    logger,
		api:       api,
		app:       app,
		publicApi: app,
	}, nil
}

func (svr *TestServer) PrivateEndpoint() fiber.Router {
	return svr.api
}

func (svr *TestServer) PublicEndpoint() fiber.Router {
	return svr.publicApi
}

func (svr *TestServer) Start(ctx context.Context) error {
	svr.logger.Sugar().Infoln("Start server")

	return svr.app.Listen(":9010")
}

func (svr *TestServer) Stop(ctx context.Context) error {
	svr.logger.Sugar().Infoln("Shutdown server")

	return svr.app.Shutdown()
}

func (svr *TestServer) Server() *fiber.App {
	return svr.app
}
