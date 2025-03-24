package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

const (
	getRoomsRequstKey = "getRoomsReq"
)

func GetRoomsSchema(c *fiber.Ctx) (interface{}, string, error) {
	qStatus := strings.Split(c.Query("status", "occupied,available,booked"), ",")
	status := make([]types.RoomStatus, len(qStatus))

	for i, s := range qStatus {
		status[i] = types.RoomStatus(s)
	}

	return types.GetRoomsRequest{
		Status:                   status,
		MongoPaginateWithIDQuery: types.NewMongoPaginateWithIDQuery(c.Query("lastID", "")),
	}, getRoomsRequstKey, nil
}
