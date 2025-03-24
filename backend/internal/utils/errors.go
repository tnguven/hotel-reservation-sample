package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tnguven/hotel-reservation-app/internal/types"
)

func ValidatorError(err error) *types.Error {
	errors := make(map[string]interface{})

	for _, v := range err.(validator.ValidationErrors) {
		suffix := ""

		if v.Tag() != "required" {
			suffix = " - invalid"
		}

		errors[v.Field()] = fmt.Sprintf("%v%v", v.Tag(), suffix)
	}

	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusBadRequest,
			Msg:    http.StatusText(http.StatusBadRequest),
			Errors: errors,
		},
	}
}

func AccessForbiddenError() *types.Error {
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusForbidden,
			Msg:    http.StatusText(http.StatusForbidden),
		},
	}
}

func NotFoundError() *types.Error {
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusNotFound,
			Msg:    http.StatusText(http.StatusNotFound),
		},
	}
}

func UnauthorizedError() *types.Error {
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusUnauthorized,
			Msg:    http.StatusText(http.StatusUnauthorized),
		},
	}
}

func ConflictError(errorMessage string) *types.Error {
	msg := http.StatusText(http.StatusConflict)
	if errorMessage != "" {
		msg = errorMessage
	}
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusConflict,
			Msg:    msg,
		},
	}
}

func InvalidCredError() *types.Error {
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusBadRequest,
			Msg:    "invalid credentials",
		},
	}
}

func BadRequestError(errorMessage string) *types.Error {
	msg := http.StatusText(http.StatusBadRequest)
	if errorMessage != "" {
		msg = errorMessage
	}
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusBadRequest,
			Msg:    msg,
		},
	}
}

func InternalServerError(errorMessage string) *types.Error {
	msg := http.StatusText(http.StatusInternalServerError)
	if errorMessage != "" {
		msg = errorMessage
	}
	return &types.Error{
		ResGeneric: &types.ResGeneric{
			Status: http.StatusInternalServerError,
			Msg:    msg,
		},
	}
}
