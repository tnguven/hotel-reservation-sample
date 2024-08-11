package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleGetUser(c *fiber.Ctx) error {
	var id = c.Params("id")
	req := getUserRequest{
		ID: id,
	}
	if err := req.bind(h.validator); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

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
	fmt.Println(c.UserContext().Value("userID"))
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

	return c.JSON(insertedUser)
}

func (h *Handler) HandleDeleteUser(c *fiber.Ctx) error {
	var id = c.Params("id")
	req := getUserRequest{
		ID: id,
	}
	if err := req.bind(h.validator); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"deleted": id})
}

func (h *Handler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	req := updateUserRequest{
		ID:        id,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}
	if err := req.bind(h.validator); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewValidatorError(err))
	}

	if err := h.userStore.PutUser(c.Context(), params, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"updated": id})
}
