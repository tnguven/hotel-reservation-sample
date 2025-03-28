package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

const (
	getHotelRequestKey  = "getHotelReq"
	getHotelsRequestKey = "getHotelsReq"
)

type getHotelRequest struct {
	HotelID string `validate:"required,id"`
}

func GetHotelRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	return &getHotelRequest{
		HotelID: c.Params("hotelID"),
	}, getHotelRequestKey, nil
}

func GetHotelsQueryRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	return &types.GetHotelsRequest{
		Rating:               c.QueryInt("rating", 0),
		Rooms:                c.QueryBool("rooms", false),
		QueryNumericPaginate: types.NewQueryNumericPaginate(c.QueryInt("limit", 10), c.QueryInt("page", 1)),
	}, getHotelsRequestKey, nil
}
