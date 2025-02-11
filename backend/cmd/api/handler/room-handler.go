package handler

import (
	"context"
	"fmt"

	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) HandleBookRoom(c *fiber.Ctx) error {
	roomID := c.Params("roomID")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return utils.UnauthorizedError()
	}

	var params types.BookingParam
	if err := c.BodyParser(&params); err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error request body")
	}
	params.RoomID = roomID
	params.UserID = user.ID

	roomIsAvailable, err := h.isBookAvailableForBooking(c.Context(), &params)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error checking if book is available")
	}
	if !roomIsAvailable {
		return c.Status(fiber.StatusConflict).JSON(&types.GenericResponse{
			Msg:    fmt.Sprintf("room %s already booked", roomID),
			Status: fiber.StatusConflict,
		})
	}

	insertedBooking, err := h.bookingStore.InsertBooking(c.Context(), &params)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error inserting booking")
	}

	return c.Status(fiber.StatusCreated).JSON(&types.GenericResponse{
		Data:   insertedBooking,
		Status: fiber.StatusCreated,
	})
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
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting rooms")
	}

	return c.Status(fiber.StatusOK).JSON(&types.GenericResponse{
		Data:   &rooms,
		Status: fiber.StatusOK,
	})
}
