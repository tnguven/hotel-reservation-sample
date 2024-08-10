package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(withLog bool) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(ErrorResponse{Error: err.Error()})
		},
	})

	if withLog {
		app.Use(etag.New())
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${status} - ${method} ${path}​\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type/application/json, Accept, ,",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	return app
}
