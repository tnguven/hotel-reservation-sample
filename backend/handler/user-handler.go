package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetUser(c *fiber.Ctx) error {
	id, _ := c.Locals("userID").(string)
	user, err := h.userStore.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		return err
	}

	return c.JSON(users)
}

func (h *Handler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exist"})
		}
		return err
	}

	return c.JSON(insertedUser)
}

func (h *Handler) HandleDeleteUser(c *fiber.Ctx) error {
	id, _ := c.Locals("userID").(string)

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"deleted": id})
}

func (h *Handler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id, _  = c.Locals("userID").(string)
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.userStore.PutUser(c.Context(), params, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"updated": id})
}
