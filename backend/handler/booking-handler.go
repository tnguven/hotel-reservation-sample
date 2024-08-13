package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// TODO: this needs to be admin authorized
func (h *Handler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.GetBookings(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusFound).JSON(bookings)
}

// TODO: this needs to be user authorized
func (h *Handler) HandleGetBooking(c *fiber.Ctx) error {
	bookingID := c.Params("bookingID")

	booking, err := h.bookingStore.GetBookingsByID(c.Context(), bookingID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusFound).JSON(booking)
}

func (h *Handler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("bookingID")
	user, err := getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if user.IsAdmin {
		if err := h.bookingStore.CancelBookingByAdmin(c.Context(), bookingID); err != nil {
			return err
		}
	} else {
		if err := h.bookingStore.CancelBookingByUserID(c.Context(), bookingID, user.ID); err != nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(genericResp{
		Type: "success",
		Msg:  fmt.Sprintf("booking %s has been canceled", bookingID),
	})
}
