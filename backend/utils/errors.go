package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	GenericResponse
}

func (err Error) Error() string {
	return err.Msg
}

func NewError(err error, code int, msg string) Error {
	errors := make(map[string]interface{})
	// add switch other variant
	switch v := err.(type) {
	default:
		errors["body"] = v.Error()
	}

	message := http.StatusText(code)
	if msg != "" {
		message = msg
	}

	return Error{
		GenericResponse: GenericResponse{
			Status: code,
			Msg:    message,
			Errors: errors,
		},
	}
}

func ValidatorError(err error) Error {
	errors := make(map[string]interface{})

	for _, v := range err.(validator.ValidationErrors) {
		suffix := ""

		if v.Tag() != "required" {
			suffix = " - invalid"
		}

		errors[v.Field()] = fmt.Sprintf("%v%v", v.Tag(), suffix)
	}

	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusBadRequest,
			Msg:    http.StatusText(http.StatusBadRequest),
			Errors: errors,
		},
	}
}

func AccessForbiddenError() Error {
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusForbidden,
			Msg:    http.StatusText(http.StatusForbidden),
		},
	}
}

func NotFoundError() Error {
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusNotFound,
			Msg:    http.StatusText(http.StatusNotFound),
		},
	}
}

func UnauthorizedError() Error {
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusUnauthorized,
			Msg:    http.StatusText(http.StatusUnauthorized),
		},
	}
}

func ConflictError(errorMessage string) Error {
	msg := http.StatusText(http.StatusConflict)
	if errorMessage != "" {
		msg = errorMessage
	}
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusConflict,
			Msg:    msg,
		},
	}
}

func InvalidCredError() Error {
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusBadRequest,
			Msg:    "invalid credentials",
		},
	}
}

func BadRequestError() Error {
	return Error{
		GenericResponse: GenericResponse{
			Status: http.StatusBadRequest,
			Msg:    http.StatusText(http.StatusBadRequest),
		},
	}
}
