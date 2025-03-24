package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"github.com/tnguven/hotel-reservation-app/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetHotels(c *fiber.Ctx) error {
	qParams, ok := c.Locals(getHotelsRequestKey).(*types.GetHotelsRequest)
	if !ok {
		log.Errorf("locals %s field missing", getHotelRequestKey)
		return utils.BadRequestError("")
	}

	hotels, total, err := h.hotelStore.GetHotels(c.Context(), qParams)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting hotels")
	}

	return c.Status(fiber.StatusOK).JSON(types.ResWithPaginate[types.ResNumericPaginate]{
		ResGeneric: types.ResGeneric{
			Data:   &hotels,
			Status: fiber.StatusOK,
		},
		Pagination: types.ResNumericPaginate{
			Count: total,
			Page:  qParams.Page,
			Limit: int(qParams.Limit),
		},
	})
}

func (h *Handler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	hotelID := c.Params("hotelID")

	rooms, err := h.roomStore.GetRoomsByHotelID(c.Context(), hotelID)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting rooms")
	}

	return c.Status(fiber.StatusOK).JSON(&types.ResGeneric{
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

	return c.Status(fiber.StatusFound).JSON(&types.ResGeneric{
		Data:   hotel,
		Status: fiber.StatusFound,
	})
}
