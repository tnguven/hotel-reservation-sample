package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	req := authRequest{
		Email:    authParams.Email,
		Password: authParams.Password,
	}
	if err := req.bind(h.validator); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credential")
		}
		return err
	}

	if !authParams.IsValidPassword(user.EncryptedPassword) {
		return fmt.Errorf("invalid credentials")
	}

	token := utils.GenerateJWT(user.ID.Hex())

	return c.JSON(&AuthResponse{
		User:  user,
		Token: token,
	})
	// return c.Status(fiber.StatusFound).JSON(user)
}

func (h *Handler) HandleSignIn(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	req := insertUserRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}
	if err := req.bind(h.validator); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
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
