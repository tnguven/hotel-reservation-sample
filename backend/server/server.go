package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tnguven/hotel-reservation-app/utils"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(withLog bool) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if utilsError, ok := err.(utils.Error); ok {
				return c.Status(utilsError.Code).JSON(utilsError)
			}

			return c.Status(fiber.StatusInternalServerError).
				JSON(utils.NewError(err, fiber.StatusInternalServerError))
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
