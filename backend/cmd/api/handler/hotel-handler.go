package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetHotels(c *fiber.Ctx) error {
	qParams, ok := c.Locals(getHotelsRequestKey).(*types.GetHotelsRequest)
	if !ok {
		log.Errorf("locals %s field missing", getHotelRequestKey)
		return utils.BadRequestError("")
	}
	fmt.Printf("qPARAMS: %+v", qParams)

	hotels, total, err := h.hotelStore.GetHotels(c.Context(), qParams)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting hotels")
	}

	return c.Status(fiber.StatusOK).JSON(&types.GenericResponse{
		Data:   &hotels,
		Status: fiber.StatusOK,
		PaginationResponse: &types.PaginationResponse{
			Count: total,
			Page:  qParams.Page,
			Limit: qParams.Limit,
		},
	})
}

func (h *Handler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")

	rooms, err := h.roomStore.GetRoomsByHotelID(c.Context(), hotelID)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting rooms")
	}

	return c.Status(fiber.StatusOK).JSON(&types.GenericResponse{
		Data:   &rooms,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), hotelID)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting hotel")
	}

	return c.Status(fiber.StatusFound).JSON(&types.GenericResponse{
		Data:   hotel,
		Status: fiber.StatusFound,
	})
}
