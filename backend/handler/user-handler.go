package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetUser(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	req := getUserRequest{
		ID: id,
	}
	err := req.bind(h.validator)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	user, err := h.userStore.GetByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.NotFound())
		}
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		return err
	}

	return ctx.JSON(users)
}

func (h *Handler) HandlePostUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	req := insertUserRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}
	err := req.bind(h.validator)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(ctx.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exist"})
		}
		return err
	}

	return ctx.JSON(insertedUser)
}

func (h *Handler) HandleDeleteUser(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	req := getUserRequest{
		ID: id,
	}
	err := req.bind(h.validator)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	if err := h.userStore.DeleteUser(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"deleted": id})
}

func (h *Handler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		id     = ctx.Params("id")
		params *types.UpdateUserParams
	)

	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	req := updateUserRequest{
		ID:        id,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}
	err := req.bind(h.validator)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	if err := h.userStore.PutUser(ctx.Context(), params, id); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"updated": id})
}
