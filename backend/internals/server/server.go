package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tnguven/hotel-reservation-app/internals/configure"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(configs configure.Common) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if response, ok := err.(*types.Error); ok {
				return c.Status(response.Status).JSON(&response)
			}

			return c.Status(fiber.StatusInternalServerError).
				JSON(types.NewError(err, fiber.StatusInternalServerError, ""))
		},
	})

	if configs.WithLog() {
		app.Use(etag.New())
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${status} - ${method} ${path}\n",
		}))
	}

	if configs.GoEnv() != "production" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type/application/json, Accept, ,",
			AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
		}))
	}

	return app
}
