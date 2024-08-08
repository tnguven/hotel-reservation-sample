package handler

import (
	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) Register(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Get("/ping", handler.HandlerPing)

	publicUser := v1.Group("/users")
	publicUser.Get("/", handler.HandleGetUsers)
	publicUser.Get("/:id", handler.HandleGetUser)
	publicUser.Put("/:id", handler.HandlePutUser)
	publicUser.Delete("/:id", handler.HandleDeleteUser)
	publicUser.Post("/", handler.HandlePostUser)

	// jwtMiddleware := jwtware.New(
	// 	jwtware.Config{
	// 		SigningKey: utils.JWTSecret,
	// 		AuthScheme: "Token",
	// 	})

	// user := v1.Group("/user", jwtMiddleware)
	// user.Get("", h.CurrentUser)
	// user.Put("", h.UpdateUser)
}
