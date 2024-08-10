package handler

import (
	"github.com/tnguven/hotel-reservation-app/store"
)

type Handler struct {
	validator  *Validator
	userStore  store.UserStore
	hotelStore store.HotelStore
	roomStore  store.RoomStore
}

func NewHandler(stores *store.Stores) *Handler {
	v := NewValidator()

	return &Handler{
		validator:  v,
		userStore:  stores.User,
		hotelStore: stores.Hotel,
		roomStore:  stores.Room,
	}
}
