package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Errors  map[string]interface{} `json:"errors"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

// add switch other variant
func NewError(err error, code int) Error {
	errors := make(map[string]interface{})

	switch v := err.(type) {
	default:
		errors["body"] = v.Error()
	}

	return Error{
		Code:    code,
		Message: http.StatusText(code),
		Errors:  errors,
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
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Errors:  errors,
	}
}

func AccessForbiddenError() Error {
	return Error{
		Message: http.StatusText(http.StatusForbidden),
		Code:    http.StatusForbidden,
	}
}

func NotFoundError() Error {
	return Error{
		Message: http.StatusText(http.StatusNotFound),
		Code:    http.StatusNotFound,
	}
}

func UnauthorizeError() Error {
	return Error{
		Message: http.StatusText(http.StatusNotFound),
		Code:    http.StatusNotFound,
	}
}
