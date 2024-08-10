package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/handler/middleware"
)

func (handler *Handler) Register(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	// auth handlers
	auth := v1.Group("/auth")
	auth.Post("/", handler.HandleAuthenticate)
	auth.Post("/signin", handler.HandleSignIn)

	// user handlers
	usersPrivate := v1.Group("/users", middleware.JWTAuthentication)
	usersPrivate.Get("/", handler.HandleGetUsers)
	usersPrivate.Get("/:id", handler.HandleGetUser)
	usersPrivate.Put("/:id", handler.HandlePutUser)
	usersPrivate.Delete("/:id", handler.HandleDeleteUser)
	usersPrivate.Post("/", handler.HandlePostUser)

	// hotel handlers
	hotels := v1.Group("/hotels")
	hotels.Get("/", handler.HandleGetHotels)
	hotels.Get("/:hotelID", handler.HandleGetHotel)
	hotels.Get("/:hotelID/rooms", handler.HandleGetRooms)

	// jwtMiddleware := jwtware.New(
	// 	jwtware.Config{
	// 		SigningKey: utils.JWTSecret,
	// 		AuthScheme: "Token",
	// 	})

	// user := v1.Group("/user", jwtMiddleware)
	// user.Get("", h.CurrentUser)
	// user.Put("", h.UpdateUser)
}
