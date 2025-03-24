package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internal/tokener"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthResponse struct {
		User  *types.User `json:"user"`
		Token string      `json:"token"`
	}

	signInRequest struct {
		*types.CreateUserParams
	}
)

func (h *Handler) HandleAuthenticate(configs tokener.JWTConfigs) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authReqParams, ok := c.Locals(authRequestKey).(*authRequest)
		if !ok {
			log.Printf("something went wrong to get the authReqParams")
			return utils.BadRequestError("")
		}
		params := types.AuthParams{
			Email:    authReqParams.Email,
			Password: authReqParams.Password,
		}

		user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return utils.InvalidCredError()
			}
			return types.NewError(err, fiber.StatusInternalServerError, "")
		}

		if !params.IsValidPassword(user.EncryptedPassword) {
			return utils.InvalidCredError()
		}

		token, err := tokener.GenerateJWT(user.ID.Hex(), user.IsAdmin, configs)
		if err != nil {
			return types.NewError(err, fiber.StatusInternalServerError, "")
		}

		return c.Status(fiber.StatusOK).JSON(&types.ResGeneric{
			Data: &AuthResponse{
				User:  user,
				Token: token,
			},
			Status: fiber.StatusOK,
		})
	}
}

func (h *Handler) HandleSignIn(c *fiber.Ctx) error {
	params, ok := c.Locals(insertUserRequestKey).(*types.CreateUserParams)
	if !ok {
		log.Println("insertUserRequest local missing")
		return utils.BadRequestError("")
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Printf("new user from params failed: %v", err)
		return utils.InternalServerError("")
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.ConflictError("email already exist")
		}

		log.Printf("userStore.InserUser failed: %v", err)
		return utils.InternalServerError("can not inset the new user...")
	}

	return c.Status(fiber.StatusCreated).JSON(&types.ResGeneric{
		Data:   insertedUser,
		Status: fiber.StatusCreated,
	})
}
