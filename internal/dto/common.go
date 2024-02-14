package dto

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Message string `json:"message,omitempty"`
}

type DownloadSignedURL struct {
	URL string `json:"url"`
}

type Error struct {
	MessageResponse
	Err        error `json:"-"`
	HttpStatus int   `json:"-"`
}

type MessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	DefaultSuccessMessage = "success"
	DefaultSuccessCode    = 0
)

type Payload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func WriteJSON(c *fiber.Ctx, payload Payload) error {
	return c.Status(http.StatusOK).
		JSON(payload)
}
