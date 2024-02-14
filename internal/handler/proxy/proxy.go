package proxy

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/internal/dto"
	"github.com/pitsanujiw/go-boilerplate/pkg/server"
)

type Handler interface {
	RegisterRouters()
}

type handler struct {
	server server.HTTPServer
	cfg    *config.App
}

func New(server server.HTTPServer, cfg *config.App) Handler {
	h := &handler{
		server: server,
		cfg:    cfg,
	}

	h.RegisterRouters()

	return h
}

func (h *handler) RegisterRouters() {
	router := h.server.PublicEndpoint()
	router.All("/endpoint/*", func(c *fiber.Ctx) error {
		return h.reverseProxy(c, h.cfg.Gateway.Endpoint)
	})
}

func (h *handler) reverseProxy(c *fiber.Ctx, endpoint string) error {
	path, ok := c.AllParams()["*1"]
	if !ok {
		return dto.WriteJSON(c, dto.Payload{
			Code:    502,
			Message: "Bad Gateway",
		})
	}

	url := fmt.Sprintf("%s/%s", endpoint, path)
	if err := proxy.Do(c, url); err != nil {
		return dto.WriteJSON(c, dto.Payload{
			Code:    502,
			Message: "Bad Gateway",
		})
	}

	return nil
}
