package handler

import (
	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) Register(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	users := v1.Group("/users")
	users.Get("/", handler.HandleGetUsers)
	users.Get("/:id", handler.HandleGetUser)
	users.Put("/:id", handler.HandlePutUser)
	users.Delete("/:id", handler.HandleDeleteUser)
	users.Post("/", handler.HandlePostUser)

	hotels := v1.Group("/hotels")
	hotels.Get("/", handler.HandleGetHotels)
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
