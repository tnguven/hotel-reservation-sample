package handler

import (
	"github.com/gofiber/fiber/v2"
	mid "github.com/tnguven/hotel-reservation-app/handler/middleware"
)

func (handler *Handler) Register(app *fiber.App) {
	v := mid.NewValidator()

	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	auth := v1.Group("/auth")
	auth.Post("/", mid.WithValidation(v, AuthRequestSchema), handler.HandleAuthenticate)

	auth.Post("/signin", mid.WithValidation(v, AuthRequestSchema), handler.HandleSignIn)

	usersPrivate := v1.Group("/users", mid.JWTAuthentication(handler.userStore))
	usersPrivate.Get("/", handler.HandleGetUsers)
	usersPrivate.Post("/", mid.WithValidation(v, InsertUserRequestSchema), handler.HandlePostUser)
	userPrivate := usersPrivate.Group("/:id")
	userPrivate.Get("/", handler.HandleGetUser)
	userPrivate.Put("/", handler.HandlePutUser)
	userPrivate.Delete("/", handler.HandleDeleteUser)

	hotelsPrivate := v1.Group("/hotels", mid.JWTAuthentication(handler.userStore))
	hotelsPrivate.Get("/", handler.HandleGetHotels)
	hotelPrivate := hotelsPrivate.Group("/:hotelID", mid.WithValidation(v, GetHotelRequestSchema))
	hotelPrivate.Get("/", handler.HandleGetHotel)
	hotelPrivate.Get("/rooms", handler.HandleGetRoomsByHotelID)

	roomsPrivate := v1.Group("/rooms", mid.JWTAuthentication(handler.userStore))
	roomsPrivate.Get("/", handler.HandleGetRooms)
	bookPrivate := roomsPrivate.Group("/:roomID") // TODO: add roomID validation
	bookPrivate.Post("/book", mid.WithValidation(v, BookingRoomRequestSchema), handler.HandleBookRoom)

	bookingsPrivate := v1.Group("/bookings", mid.JWTAuthentication(handler.userStore))
	bookingsPrivate.Get("/", mid.WithAdminAuth, handler.HandleGetBookings)
	bookingsPrivate.Get("/:bookingID", handler.HandleGetBooking)           // TODO: validate id
	bookingsPrivate.Put("/:bookingID/cancel", handler.HandleCancelBooking) // TODO: validate id
}
