package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetHotels(c *fiber.Ctx) error {
	var qParams store.HotelQueryParams
	if err := c.QueryParser(&qParams); err != nil {
		return err
	}
	fmt.Println(qParams) // TODO: implement

	user, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")

	rooms, err := h.roomStore.GetRoomsByHotelID(c.Context(), hotelID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}

func (h *Handler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")

	hotel, err := h.hotelStore.GetHotelByID(c.Context(), hotelID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusFound).JSON(hotel)
}
