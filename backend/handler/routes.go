package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/config"
	mid "github.com/tnguven/hotel-reservation-app/middleware"
)

func (h *Handler) Register(app *fiber.App, configs *config.Configs) {
	v := mid.NewValidator()

	v1 := app.Group("/v1")
	v1.Get("/ping", h.HandlerPing)

	withAuth := mid.JWTAuthentication(h.userStore, configs)

	auth := v1.Group("/auth")
	auth.Post("/", mid.WithValidation(v, AuthRequestSchema), h.HandleAuthenticate(configs))
	auth.Post("/signin", mid.WithValidation(v, InsertUserRequestSchema), h.HandleSignIn)

	usersPrivate := v1.Group("/users")
	usersPrivate.Get("/", withAuth, h.HandleGetUsers)
	usersPrivate.Post("/", mid.WithValidation(v, InsertUserRequestSchema), h.HandlePostUser)

	userPrivate := usersPrivate.Group("/:id")
	userPrivate.Get("/", h.HandleGetUser)
	userPrivate.Put("/", mid.WithValidation(v, UpdateUserRequestSchema), h.HandlePutUser)
	userPrivate.Delete("/", h.HandleDeleteUser)

	hotelsPrivate := v1.Group("/hotels", withAuth)
	hotelsPrivate.Get("/", h.HandleGetHotels)
	hotelPrivate := hotelsPrivate.Group("/:hotelID", mid.WithValidation(v, GetHotelRequestSchema))
	hotelPrivate.Get("/", h.HandleGetHotel)
	hotelPrivate.Get("/rooms", h.HandleGetRoomsByHotelID)

	roomsPrivate := v1.Group("/rooms", withAuth)
	roomsPrivate.Get("/", h.HandleGetRooms)
	bookPrivate := roomsPrivate.Group("/:roomID") // TODO: add roomID validation
	bookPrivate.Post("/book", mid.WithValidation(v, BookingRoomRequestSchema), h.HandleBookRoom)

	adminBookings := v1.Group("/admin/bookings", withAuth)
	adminBookings.Get("/", mid.WithAdminAuth, h.HandleGetBookingsAsAdmin)

	bookingsPrivate := v1.Group("/bookings", withAuth)
	bookingsPrivate.Get("/", h.HandleGetBookingsAsUser)
	bookingsPrivate.Get("/:bookingID", h.HandleGetBooking)           // TODO: validate id
	bookingsPrivate.Put("/:bookingID/cancel", h.HandleCancelBooking) // TODO: validate id
}
