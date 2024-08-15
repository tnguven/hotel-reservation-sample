package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredResp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credential",
	})
}

func (h *Handler) HandleAuthenticate(configs *config.Configs) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var authParams types.AuthParams
		if err := c.BodyParser(&authParams); err != nil {
			return err
		}

		user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return invalidCredResp(c)
			}
			return err
		}

		if !authParams.IsValidPassword(user.EncryptedPassword) {
			return invalidCredResp(c)
		}

		token := utils.GenerateJWT(user.ID.Hex(), user.IsAdmin, configs)

		return c.JSON(&AuthResponse{
			User:  user,
			Token: token,
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
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exist"})
		}
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(insertedUser)
}
