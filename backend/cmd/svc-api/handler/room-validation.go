package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

const (
	getRoomsRequstKey = "getRoomsReq"
)

func GetRoomsSchema(c *fiber.Ctx) (interface{}, string, error) {
	return &types.GetRoomsRequest{
		Status:          []types.RoomStatus{"occupied"},
		PaginationQuery: types.NewPaginationQuery(c.QueryInt("limit", 10), c.QueryInt("page", 1)),
	}, getRoomsRequstKey, nil
}
