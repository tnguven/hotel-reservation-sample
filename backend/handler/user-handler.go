package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		log.Errorf("HandleGetUser: error getting user: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "Error getting user")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   user,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		log.Errorf("HandleGetUsers: error getting users: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "error getting users")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Data:   &users,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return utils.NewError(err, fiber.StatusInternalServerError, "")
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(&utils.GenericResponse{
				Msg:    "email already exist",
				Status: fiber.StatusConflict,
			})
		}
		log.Errorf("HandlePostUser: error inserting user: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "something went wrong")
	}

	return c.Status(fiber.StatusCreated).JSON(&utils.GenericResponse{
		Data:   insertedUser,
		Status: fiber.StatusCreated,
	})
}

func (h *Handler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		log.Errorf("HandleDeleteUser: error deleting user: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "error deleting user")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Msg:    fmt.Sprintf("User %s deleted", id),
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)
	if err := c.BodyParser(&params); err != nil {
		log.Errorf("HandlePutUser: error parsing params: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "error parsing body")
	}
	if err := h.userStore.PutUser(c.Context(), params, id); err != nil {
		log.Errorf("HandlePutUser: error putting user: %v", err)
		return utils.NewError(err, fiber.StatusInternalServerError, "error updating user")
	}

	return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
		Msg:    fmt.Sprintf("User %v updated", id),
		Status: fiber.StatusOK,
	})
}
