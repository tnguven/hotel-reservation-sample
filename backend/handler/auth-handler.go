package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleAuthenticate(configs *config.Configs) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var authParams types.AuthParams
		if err := c.BodyParser(&authParams); err != nil {
			return err
		}

		user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return utils.InvalidCredError()
			}
			return utils.NewError(err, fiber.StatusInternalServerError, "")
		}

		if !authParams.IsValidPassword(user.EncryptedPassword) {
			return utils.InvalidCredError()
		}

		token, err := utils.GenerateJWT(user.ID.Hex(), user.IsAdmin, configs)
		if err != nil {
			return utils.NewError(err, fiber.StatusInternalServerError, "")
		}

		return c.Status(fiber.StatusOK).JSON(&utils.GenericResponse{
			Data: &AuthResponse{
				User:  user,
				Token: token,
			},
			Status: fiber.StatusOK,
		})
	}
}

func (h *Handler) HandleSignIn(c *fiber.Ctx) error {
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
			return utils.ConflictError("email already exist")
		}
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&utils.GenericResponse{
		Data:   insertedUser,
		Status: fiber.StatusCreated,
	})
}
