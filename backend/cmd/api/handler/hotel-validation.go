package handler

import (
	"github.com/gofiber/fiber/v2"
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

type hotelQueryRequest struct {
	Rooms  bool `validate:"boolean,omitempty"`
	Rating int  `validate:"numeric,omitempty"`
	*PaginationFilter
}

func GetHotelsQueryRequestSchema(c *fiber.Ctx) (interface{}, string, error) {
	return &hotelQueryRequest{
		Rating: c.QueryInt("rating", 0),
		Rooms:  c.QueryBool("rooms", false),
		PaginationFilter: &PaginationFilter{
			Limit: int64(c.QueryInt("limit", 10)), // TODO: manage the hardcoded defaults
			Page:  int64(c.QueryInt("page", 1)),
		},
	}, getHotelsRequestKey, nil
}
