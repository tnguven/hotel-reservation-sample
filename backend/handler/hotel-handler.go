package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetHotels(c *fiber.Ctx) error {
	var qParams store.HotelQueryParams
	if err := c.QueryParser(&qParams); err != nil {
		log.Error("Error parsing query parameters", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "Error query parameters")
	}

	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting hotels")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   &hotels,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")

	rooms, err := h.roomStore.GetRoomsByHotelID(c.Context(), hotelID)
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting rooms")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   &rooms,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), hotelID)
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting hotel")
	}

	return c.Status(fiber.StatusFound).JSON(&utils.GenericResponse{
		Data:   hotel,
		Status: fiber.StatusFound,
	})
}
