package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

func (h *Handler) HandleGetBookingsAsUser(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return utils.UnauthorizedError()
	}

	bookings, err := h.bookingStore.GetBookingsAsUser(c.Context(), user)
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting bookings")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   &bookings,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetBookingsAsAdmin(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.GetBookingsAsAdmin(c.Context())
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting bookings")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   &bookings,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetBooking(c *fiber.Ctx) error {
	bookingID := c.Params("bookingID")

	booking, err := h.bookingStore.GetBookingsByID(c.Context(), bookingID)
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting booking")
	}

	return c.Status(fiber.StatusFound).JSON(&utils.GenericResponse{
		Data:   booking,
		Status: fiber.StatusFound,
	})
}

func (h *Handler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("bookingID")
	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return utils.UnauthorizedError()
	}

	if user.IsAdmin {
		if err = h.bookingStore.CancelBookingByAdmin(c.Context(), bookingID); err != nil {
			return utils.NewError(err, fiber.StatusInternalServerError, "Failed to cancel booking")
		}
	} else {
		if err = h.bookingStore.CancelBookingByUserID(c.Context(), bookingID, user.ID); err != nil {
			return utils.NewError(err, fiber.StatusInternalServerError, "Failed to cancel booking")
		}
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Msg:    fmt.Sprintf("booking %s has been canceled", bookingID),
		Status: fiber.StatusOK,
	})
}
