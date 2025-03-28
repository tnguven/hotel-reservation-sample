package handler

import (
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

	insertedBooking, err := h.bookingStore.InsertBooking(c.Context(), &params)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error inserting booking")
	}

	return c.Status(fiber.StatusCreated).JSON(types.ResGeneric{
		Data:   insertedBooking,
		Status: fiber.StatusCreated,
	})
}

func (h *Handler) HandleGetRooms(c *fiber.Ctx) error {
	qParams := c.Locals(getRoomsRequstKey).(*types.GetRoomsRequest)
	rooms, total, nextLastId, err := h.roomStore.GetRooms(c.Context(), qParams)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting rooms")
	}

	return c.Status(fiber.StatusOK).JSON(types.ResWithPaginate[types.ResCursorPaginate]{
		ResGeneric: types.ResGeneric{
			Data:   rooms,
			Status: fiber.StatusOK,
		},
		Pagination: types.ResCursorPaginate{
			LastID: nextLastId,
			Limit:  int(qParams.Limit),
			Count:  total,
		},
	})
}
