package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
)

type GetUserRequest struct {
	ID string `validate:"required,id"`
}

func GetUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	id := c.Params("id")
	return &GetUserRequest{
		ID: id,
	}, nil
}

type insertUserRequest struct {
	FirstName string `validate:"required,alpha,min=2,max=48"`
	LastName  string `validate:"required,alpha,min=2,max=48"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=7,max=256"`
}

func InsertUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return nil, err
	}

	return &insertUserRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}, nil
}

type updateUserRequest struct {
	ID        string `validate:"required,id"`
	FirstName string `validate:"omitempty,alpha,min=2,max=48"`
	LastName  string `validate:"omitempty,alpha,min=2,max=48"`
}

func UpdateUserRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return nil, err
	}
	return &updateUserRequest{
		ID:        id,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, nil
}

type getHotelRequest struct {
	HotelID string `validate:"required,id"`
}

func GetHotelRequestSchema(c *fiber.Ctx) (interface{}, error) {
	hotelID := c.Params("hotelID")
	return &getHotelRequest{
		HotelID: hotelID,
	}, nil
}

type hotelQueryRequest struct {
	Rooms  bool `validate:"rooms"`
	Rating int  `validate:"numeric"`
}

func HotelQueryRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var queryParams store.HotelQueryParams
	if err := c.QueryParser(&queryParams); err != nil {
		return nil, err
	}

	return &hotelQueryRequest{
		Rooms:  queryParams.Rooms,
		Rating: queryParams.Rating,
	}, nil
}

type authRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=7,max=256"`
}

func AuthRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return nil, err
	}

	return &authRequest{
		Email:    authParams.Email,
		Password: authParams.Password,
	}, nil
}

type bookingRoomRequest struct {
	FromDate  time.Time `validate:"required"`
	TillDate  time.Time `validate:"required"`
	NumPerson int       `validate:"required,numeric,min=1,max=20"`
}

func BookingRoomRequestSchema(c *fiber.Ctx) (interface{}, error) {
	var bookingParams types.BookingParam
	if err := c.BodyParser(&bookingParams); err != nil {
		return nil, err
	}
	return &bookingRoomRequest{
		FromDate:  bookingParams.FromDate,
		TillDate:  bookingParams.TillDate,
		NumPerson: bookingParams.CountPerson,
	}, nil
}
