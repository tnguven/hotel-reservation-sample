package handler

import (
	"github.com/tnguven/hotel-reservation-app/store"
)

type Handler struct {
	validator *Validator
	userStore store.UserStore
}

func NewHandler(userStore store.UserStore) *Handler {
	v := NewValidator()

	return &Handler{
		validator: v,
		userStore: userStore,
	}
}
