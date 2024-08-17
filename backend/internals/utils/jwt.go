package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tnguven/hotel-reservation-app/internals/config"
)

func GenerateJWT(id string, isAdmin bool, configs *config.Configs) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"exp":     time.Now().Add(time.Hour * time.Duration(configs.TokenExpHour)).Unix(),
		"isAdmin": isAdmin,
	})
	t, err := token.SignedString([]byte(configs.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token with secret")
	}
	return t, nil
}
