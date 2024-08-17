package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

type bookingRoomRequest struct {
	FromDate  time.Time `validate:"required"`
	TillDate  time.Time `validate:"required"`
	NumPerson int       `validate:"required,numeric,min=1,max=20"`
}

func BookingRoomRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var params types.BookingParam
	if err := c.BodyParser(&params); err != nil {
		return nil, err
	}

	now := time.Now()
	if now.After(params.FromDate) || now.After(params.TillDate) {
		return nil, fmt.Errorf("cannot book a room in the past")
	}

	return &bookingRoomRequest{
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPerson: params.CountPerson,
	}, nil
}
