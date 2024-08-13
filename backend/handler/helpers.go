package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
)

func getAuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return user, nil
}

// func sendUnauthorized(c *fiber.Ctx) error {
// 	return c.Status(fiber.StatusUnauthorized).JSON(genericResp{
// 		Msg:  "not authorized",
// 		Type: "error",
// 	})
// }
