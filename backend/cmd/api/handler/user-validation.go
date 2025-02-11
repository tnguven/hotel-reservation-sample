package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

const (
	defaultReadLimit     = 10
	insertUserRequestKey = "insertUserReqKey"
	updateUserRequestKey = "updateUserReqKey"
	getUserRequestKey    = "getUserReqKey"
	getUsersRequestKey   = "getUsersReqKey"
)

func InsertUserRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return nil, insertUserRequestKey, utils.BadRequestError(err.Error())
	}

	return &params, insertUserRequestKey, nil
}

type updateUserRequest struct {
	ID        string `validate:"required,id" json:"id"`
	FirstName string `validate:"omitempty,alpha,min=2,max=48" json:"firstName"`
	LastName  string `validate:"omitempty,alpha,min=2,max=48" json:"lastName"`
}

func UpdateUserRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return nil, updateUserRequestKey, utils.BadRequestError(err.Error())
	}

	return &updateUserRequest{
		ID:        id,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, updateUserRequestKey, nil
}

type getUserRequest struct {
	ID string `validate:"required,id" json:"id"`
}

func GetUserRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	return &getUserRequest{
		ID: c.Params("id"),
	}, getUserRequestKey, nil
}

func GetUsersRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	return types.NewPaginationQuery(limit, page), getUsersRequestKey, nil
}
