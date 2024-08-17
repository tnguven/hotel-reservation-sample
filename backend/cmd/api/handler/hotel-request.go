package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/store"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

type getHotelRequest struct {
	HotelID string `validate:"required,id"`
}

func GetHotelRequestSchema(c *fiber.Ctx) (interface{}, error) {
	hotelID := c.Params("hotelID")
	return &getHotelRequest{
		HotelID: hotelID,
	}, nil
}

type hotelQueryRequest struct {
	Rating int `validate:"numeric"`
}

func HotelQueryRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var queryParams store.HotelQueryParams
	if err := c.QueryParser(&queryParams); err != nil {
		return nil, utils.BadRequestError()
	}

	return &hotelQueryRequest{
		Rating: queryParams.Rating,
	}, nil
}
