package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/config"
	mid "github.com/tnguven/hotel-reservation-app/internals/middleware"
)

func (h *Handler) Register(app *fiber.App, configs *config.Configs, validator *mid.Validator) {
	v1 := app.Group("/v1")

	withAuth := mid.JWTAuthentication(h.userStore, configs)

	auth := v1.Group("/auth")
	auth.Post("/", mid.WithValidation(validator, AuthRequestSchema), h.HandleAuthenticate(configs))
	auth.Post("/signin", mid.WithValidation(validator, InsertUserRequestSchema), h.HandleSignIn)

	usersPrivate := v1.Group("/users")
	usersPrivate.Get("/", withAuth, mid.WithValidation(validator, GetUsersRequestSchema), h.HandleGetUsers)
	usersPrivate.Post("/", mid.WithValidation(validator, InsertUserRequestSchema), h.HandlePostUser)

	userPrivate := usersPrivate.Group("/:id")
	userPrivate.Get("/", h.HandleGetUser)
	userPrivate.Put("/", mid.WithValidation(validator, UpdateUserRequestSchema), h.HandlePutUser)
	userPrivate.Delete("/", h.HandleDeleteUser)

	hotelsPrivate := v1.Group("/hotels", withAuth)
	hotelsPrivate.Get("/", h.HandleGetHotels)
	hotelPrivate := hotelsPrivate.Group("/:hotelID", mid.WithValidation(validator, GetHotelRequestSchema))
	hotelPrivate.Get("/", h.HandleGetHotel)
	hotelPrivate.Get("/rooms", h.HandleGetRoomsByHotelID)

	roomsPrivate := v1.Group("/rooms", withAuth)
	roomsPrivate.Get("/", h.HandleGetRooms)
	bookPrivate := roomsPrivate.Group("/:roomID") // TODO: add roomID validation
	bookPrivate.Post("/book", mid.WithValidation(validator, BookingRoomRequestSchema), h.HandleBookRoom)

	adminBookings := v1.Group("/admin/bookings", withAuth)
	adminBookings.Get("/", mid.WithAdminAuth, h.HandleGetBookingsAsAdmin)

	bookingsPrivate := v1.Group("/bookings", withAuth)
	bookingsPrivate.Get("/", h.HandleGetBookingsAsUser)
	bookingsPrivate.Get("/:bookingID", h.HandleGetBooking)           // TODO: validate id
	bookingsPrivate.Put("/:bookingID/cancel", h.HandleCancelBooking) // TODO: validate id
}
