package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internal/types"
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

	lastID := c.Query("lastID", "")
	limit := c.QueryInt("limit", 10)

	fmt.Println("---------", lastID)
	fmt.Println("LIMIT", limit)

	return &types.GetRoomsRequest{
		Status: status,
		QueryCursorPaginate: types.NewMongoQueryCursorPaginate(
			lastID,
			limit,
		),
	}, getRoomsRequstKey, nil
}
