package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"github.com/tnguven/hotel-reservation-app/internal/utils"

	"github.com/gofiber/fiber/v2"
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
		return types.NewError(err, fiber.StatusInternalServerError, "Error getting user")
	}

	return c.Status(fiber.StatusOK).JSON(types.ResGeneric{
		Data:   user,
		Status: fiber.StatusOK,
	})
}

func (h *Handler) HandleGetUsers(c *fiber.Ctx) error {
	query, ok := c.Locals(getUsersRequestKey).(*types.QueryNumericPaginate)
	if !ok {
		log.Error("getUsers missing locals")
		return utils.BadRequestError("")
	}

	users, total, err := h.userStore.GetUsers(c.Context(), query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return utils.NotFoundError()
		}
		log.Errorf("HandleGetUsers: error getting users: %v", err)
		return types.NewError(err, fiber.StatusInternalServerError, "error getting users")
	}

	return c.Status(fiber.StatusOK).JSON(types.ResWithPaginate[types.ResNumericPaginate]{
		ResGeneric: types.ResGeneric{
			Status: fiber.StatusOK,
			Data:   users,
		},
		Pagination: types.ResNumericPaginate{
			Count: total,
			Page:  query.Page,
			Limit: int(query.Limit),
		},
	})
}

func (h *Handler) HandlePostUser(c *fiber.Ctx) error {
	var params *types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return types.NewError(err, fiber.StatusInternalServerError, "")
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(types.ResGeneric{
				Msg:    "email already exist",
				Status: fiber.StatusConflict,
			})
		}
		log.Errorf("HandlePostUser: error inserting user: %v", err)
		return types.NewError(err, fiber.StatusInternalServerError, "something went wrong")
	}

	return c.Status(fiber.StatusCreated).JSON(types.ResGeneric{
		Data:   insertedUser,
		Status: fiber.StatusCreated,
	})
}

func (h *Handler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		log.Errorf("HandleDeleteUser: error deleting user: %v", err)
		return types.NewError(err, fiber.StatusInternalServerError, "error deleting user")
	}

	return c.Status(fiber.StatusOK).JSON(types.ResGeneric{
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
		return types.NewError(err, fiber.StatusInternalServerError, "error parsing body")
	}

	matchCount, updateErr := h.userStore.PutUser(c.Context(), params, id)
	if updateErr != nil {
		log.Errorf("HandlePutUser: error putting user: %v", updateErr)
		return types.NewError(updateErr, fiber.StatusInternalServerError, "error updating user")
	}

	if matchCount == 0 {
		return types.Error{
			ResGeneric: &types.ResGeneric{
				Status: http.StatusNotFound,
				Msg:    fmt.Sprintf("no user found with id %s", id),
			},
		}
	}

	return c.Status(fiber.StatusOK).JSON(types.ResGeneric{
		Msg:    fmt.Sprintf("User %s updated", id),
		Status: fiber.StatusOK,
	})
}
