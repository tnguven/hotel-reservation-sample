package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) HandlerPing(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{"pong": "true"})
}
