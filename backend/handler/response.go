package handler

import (
	"github.com/tnguven/hotel-reservation-app/types"
)

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
