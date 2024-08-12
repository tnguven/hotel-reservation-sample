package handler

import (
	"github.com/tnguven/hotel-reservation-app/store"
)

type Handler struct {
	userStore  store.UserStore
	hotelStore store.HotelStore
	roomStore  store.RoomStore
}

func NewHandler(stores *store.Stores) *Handler {
	return &Handler{
		userStore:  stores.User,
		hotelStore: stores.Hotel,
		roomStore:  stores.Room,
	}
}
