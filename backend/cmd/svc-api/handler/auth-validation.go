package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internal/types"
)

type authRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=7,max=256" json:"password"`
}

const authRequestKey = "authReq"

func AuthRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	var authParams *types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return nil, authRequestKey, err
	}

	return &authRequest{
		Email:    authParams.Email,
		Password: authParams.Password,
	}, authRequestKey, nil
}
