package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
)

func (h *Handler) HandleBookRoom(c *fiber.Ctx) error {
	roomID := c.Params("roomID")

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	var params types.BookingParam
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	params.RoomID = roomID
	params.UserID = user.ID
	insertedBooking, err := h.bookingStore.InsertBooking(c.Context(), &params)
	if err != nil {
		return err
	}

	fmt.Printf("booking: %+v\n", insertedBooking)

	return c.Status(fiber.StatusCreated).JSON(insertedBooking)
}
