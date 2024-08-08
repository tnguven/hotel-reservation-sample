package utils

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

var secretString = []byte("!!6359488b56d1e4196a1705f472347dd70964d8536c5e9f48e690a9eacd7b7a58!!")
var JWTSecret = jwtware.SigningKey{Key: secretString}

func GenerateJWT(id uint) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(secretString)
	return t
}
