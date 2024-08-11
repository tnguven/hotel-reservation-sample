package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/handler/middleware"
)

func (handler *Handler) Register(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	auth := v1.Group("/auth")
	auth.Post("/", handler.HandleAuthenticate)
	auth.Post("/signin", handler.HandleSignIn)

	usersPrivate := v1.Group("/users", middleware.JWTAuthentication)
	usersPrivate.Get("/", handler.HandleGetUsers)
	usersPrivate.Get("/:id", handler.HandleGetUser)
	usersPrivate.Put("/:id", handler.HandlePutUser)
	usersPrivate.Delete("/:id", handler.HandleDeleteUser)
	usersPrivate.Post("/", handler.HandlePostUser)

	hotelsPrivate := v1.Group("/hotels", middleware.JWTAuthentication)
	hotelsPrivate.Get("/", handler.HandleGetHotels)
	hotelsPrivate.Get("/:hotelID", handler.HandleGetHotel)
	hotelsPrivate.Get("/:hotelID/rooms", handler.HandleGetRooms)
}
