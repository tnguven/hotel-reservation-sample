package handler

import "github.com/gofiber/fiber/v2"

func (h *Handler) HandleHealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	}
}
