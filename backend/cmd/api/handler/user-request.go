package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

const (
	defaultReadLimit = 10
)

type insertUserRequest struct {
	FirstName string `validate:"required,alpha,min=2,max=48"`
	LastName  string `validate:"required,alpha,min=2,max=48"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=7,max=256"`
}

func InsertUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return nil, utils.BadRequestError()
	}

	return &insertUserRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}, nil
}

type updateUserRequest struct {
	ID        string `validate:"required,id"`
	FirstName string `validate:"omitempty,alpha,min=2,max=48"`
	LastName  string `validate:"omitempty,alpha,min=2,max=48"`
}

func UpdateUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return nil, utils.BadRequestError()
	}

	return &updateUserRequest{
		ID:        id,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, nil
}

type getUserRequest struct {
	ID string `validate:"required,id"`
}

func GetUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	id := c.Params("id")

	return &getUserRequest{
		ID: id,
	}, nil
}

type getUsersRequest struct {
	*PaginationFilter
}

func GetUsersRequestSchema(c *fiber.Ctx) (any, error) {
	var queryFilter PaginationFilter
	if err := c.QueryParser(&queryFilter); err != nil {
		return nil, err
	}

	return &getUsersRequest{
		PaginationFilter: &PaginationFilter{
			Limit:  queryFilter.Limit,
			Offset: queryFilter.Offset,
		},
	}, nil
}
