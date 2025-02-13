package handler

import "github.com/tnguven/hotel-reservation-app/internals/store"

type Handler struct {
	userStore    store.UserStore
	hotelStore   store.HotelStore
	roomStore    store.RoomStore
	bookingStore store.BookingStore
}

func NewHandler(stores *store.Stores) *Handler {
	return &Handler{
		userStore:    stores.User,
		hotelStore:   stores.Hotel,
		roomStore:    stores.Room,
		bookingStore: stores.Booking,
	}
}
