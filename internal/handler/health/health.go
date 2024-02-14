package health

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

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
	router.Get("/health", h.health)
}

// Health handler
//
//	@Summary	health
//	@Id			1
//	@Tags		health
//	@version	1.0
//	@produce	application/json
//	@Success	200	{object}	dto.HealthResponse
//	@Router		/health [get]
func (h *handler) health(c *fiber.Ctx) error {
	return dto.WriteJSON(c, dto.Payload{
		Message: fmt.Sprintf("health check available: %s", h.cfg.Name),
	})
}
