package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
)

func WithAdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("no permission")
	}

	if !user.IsAdmin {
		return fmt.Errorf("no permission")
	}

	return c.Next()
}
