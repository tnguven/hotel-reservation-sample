package utils

import (
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var secretString = []byte(os.Getenv("JWT_SECRET"))
var JWTSecret = jwtware.SigningKey{Key: secretString}

func GenerateJWT(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	t, err := token.SignedString(secretString)
	if err != nil {
		log.Error("failed to sign token with secret")
	}
	return t
}
