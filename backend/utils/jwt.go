package utils

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnguven/hotel-reservation-app/config"
)

func GenerateJWT(id string, isAdmin bool, configs *config.Configs) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      id,
		"exp":     time.Now().Add(time.Hour * time.Duration(configs.TokenExpHour)).Unix(),
		"isAdmin": isAdmin,
	})
	t, err := token.SignedString([]byte(configs.JWTSecret))
	if err != nil {
		log.Error("failed to sign token with secret")
	}
	return t
}
