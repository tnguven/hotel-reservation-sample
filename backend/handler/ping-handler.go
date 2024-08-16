package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/utils"
)

func (h *Handler) HandlerPing(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&utils.GenericResponse{
		Data:   fiber.Map{"pong": "pong"},
		Status: fiber.StatusOK,
	})
}
