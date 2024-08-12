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
	auth.Post("/", mid.MiddlewareValidation(v, AuthRequestSchema), handler.HandleAuthenticate)
	auth.Post("/signin", mid.MiddlewareValidation(v, AuthRequestSchema), handler.HandleSignIn)

	usersPrivate := v1.Group("/users", mid.JWTAuthentication(handler.userStore))
	usersPrivate.Get("/", handler.HandleGetUsers)
	usersPrivate.Post("/", mid.MiddlewareValidation(v, InsertUserRequestSchema), handler.HandlePostUser)
	userPrivate := usersPrivate.Group("/:id")
	userPrivate.Get("/", handler.HandleGetUser)
	userPrivate.Put("/", handler.HandlePutUser)
	userPrivate.Delete("/", handler.HandleDeleteUser)

	hotelsPrivate := v1.Group("/hotels", mid.JWTAuthentication(handler.userStore))
	hotelsPrivate.Get("/", handler.HandleGetHotels)
	hotelPrivate := hotelsPrivate.Group("/:hotelID", mid.MiddlewareValidation(v, GetHotelRequestSchema))
	hotelPrivate.Get("/", handler.HandleGetHotel)
	hotelPrivate.Get("/rooms", handler.HandleGetRooms)

	bookingPrivate := v1.Group("/room/:roomID", mid.JWTAuthentication(handler.userStore))
	bookingPrivate.Post("/book", mid.MiddlewareValidation(v, BookingRoomRequestSchema), handler.HandleBookRoom)
}
