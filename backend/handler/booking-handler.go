package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
)

func (h *Handler) HandleGetBookingsAsUser(c *fiber.Ctx) error {
	user, _ := c.Context().UserValue("user").(*types.User)
	bookings, err := h.bookingStore.GetBookingsAsUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(bookings)
}

func (h *Handler) HandleGetBookingsAsAdmin(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.GetBookingsAsAdmin(c.Context())

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(bookings)
}

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
