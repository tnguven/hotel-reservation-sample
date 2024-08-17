package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

type signInRequest struct {
	*insertUserRequest
}

func SignInRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return nil, err
	}

	return &signInRequest{
		insertUserRequest: &insertUserRequest{
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     params.Email,
			Password:  params.Password,
		},
	}, nil
}

type authRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=7,max=256"`
}

func AuthRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return nil, err
	}

	return &authRequest{
		Email:    authParams.Email,
		Password: authParams.Password,
	}, nil
}
