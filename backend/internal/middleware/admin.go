package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
)

func WithAdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok || !user.IsAdmin {
		return utils.AccessForbiddenError()
	}

	return c.Next()
}
