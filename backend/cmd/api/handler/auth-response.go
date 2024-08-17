package handler

import "github.com/tnguven/hotel-reservation-app/internals/types"

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
