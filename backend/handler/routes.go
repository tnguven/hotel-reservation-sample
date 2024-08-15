package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/config"
	mid "github.com/tnguven/hotel-reservation-app/handler/middleware"
)

func (handler *Handler) Register(app *fiber.App, configs *config.Configs) {
	v := mid.NewValidator()

	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	withAuth := mid.JWTAuthentication(handler.userStore, configs)

	auth := v1.Group("/auth")
	auth.Post("/", mid.WithValidation(v, AuthRequestSchema), handler.HandleAuthenticate(configs))
	auth.Post("/signin", mid.WithValidation(v, InsertUserRequestSchema), handler.HandleSignIn)

	usersPrivate := v1.Group("/users")
	usersPrivate.Get("/", withAuth, handler.HandleGetUsers)
	usersPrivate.Post("/", mid.WithValidation(v, InsertUserRequestSchema), handler.HandlePostUser)

	userPrivate := usersPrivate.Group("/:id")
	userPrivate.Get("/", handler.HandleGetUser)
	userPrivate.Put("/", mid.WithValidation(v, UpdateUserRequestSchema), handler.HandlePutUser)
	userPrivate.Delete("/", handler.HandleDeleteUser)

	hotelsPrivate := v1.Group("/hotels", withAuth)
	hotelsPrivate.Get("/", handler.HandleGetHotels)
	hotelPrivate := hotelsPrivate.Group("/:hotelID", mid.WithValidation(v, GetHotelRequestSchema))
	hotelPrivate.Get("/", handler.HandleGetHotel)
	hotelPrivate.Get("/rooms", handler.HandleGetRoomsByHotelID)

	roomsPrivate := v1.Group("/rooms", withAuth)
	roomsPrivate.Get("/", handler.HandleGetRooms)
	bookPrivate := roomsPrivate.Group("/:roomID") // TODO: add roomID validation
	bookPrivate.Post("/book", mid.WithValidation(v, BookingRoomRequestSchema), handler.HandleBookRoom)

	adminBookings := v1.Group("/admin/bookings", withAuth)
	adminBookings.Get("/", mid.WithAdminAuth, handler.HandleGetBookingsAsAdmin)

	bookingsPrivate := v1.Group("/bookings", withAuth)
	bookingsPrivate.Get("/", handler.HandleGetBookingsAsUser)
	bookingsPrivate.Get("/:bookingID", handler.HandleGetBooking)           // TODO: validate id
	bookingsPrivate.Put("/:bookingID/cancel", handler.HandleCancelBooking) // TODO: validate id
}
