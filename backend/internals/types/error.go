package types

import (
	"net/http"
)

type Error struct {
	*ResGeneric
}

func (err Error) Error() string {
	return err.Msg
}

func NewError(err error, code int, msg string) *Error {
	errors := make(map[string]interface{})
	// TODO: add switch other variant
	switch v := err.(type) {
	default:
		errors["body"] = v.Error()
	}

	message := http.StatusText(code)
	if msg != "" {
		message = msg
	}

	return &Error{
		ResGeneric: &ResGeneric{
			Status: code,
			Msg:    message,
			Errors: errors,
		},
	}
}
