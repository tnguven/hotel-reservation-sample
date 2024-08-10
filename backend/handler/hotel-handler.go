package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetHotels(c *fiber.Ctx) error {
	var qParams HotelQueryParams
	if err := c.QueryParser(&qParams); err != nil {
		return err
	}

	fmt.Println(qParams)

	user, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) HandleGetRooms(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")
	fmt.Println(hotelID)

	rooms, err := h.roomStore.GetRooms(c.Context(), hotelID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}
