package handler

import (
	"context"
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

	roomIsAvailable, err := h.isBookAvailableForBooking(c.Context(), &params)
	if err != nil {
		return err
	}
	if !roomIsAvailable {
		return c.Status(fiber.StatusConflict).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID),
		})
	}

	insertedBooking, err := h.bookingStore.InsertBooking(c.Context(), &params)
	if err != nil {
		return err
	}

	fmt.Printf("booking: %+v\n", insertedBooking)

	return c.Status(fiber.StatusCreated).JSON(insertedBooking)
}

func (h *Handler) isBookAvailableForBooking(ctx context.Context, params *types.BookingParam) (bool, error) {
	bookings, err := h.bookingStore.GetBookingsByRoomID(ctx, params)
	if err != nil {
		return false, err
	}

	return len(bookings) == 0, nil
}

func (h *Handler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.roomStore.GetRooms(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
