package handler

import "github.com/gofiber/fiber/v2"

func (h *Handler) HandleNotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Page not found",
	})
}
