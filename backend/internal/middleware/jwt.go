package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnguven/hotel-reservation-app/internal/configure"
	"github.com/tnguven/hotel-reservation-app/internal/store"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
)

func JWTAuthentication(userStore store.UserStore, configs configure.Secrets) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return utils.UnauthorizedError()
		}

		claims, err := validateToken(token[0], configs.JWTSecret())
		if err != nil {
			return utils.BadRequestError("missing api token")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetByID(c.Context(), userID)
		if err != nil {
			return utils.UnauthorizedError()
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid signing method", t.Header["alg"])
			return nil, utils.UnauthorizedError()
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.UnauthorizedError()
	}

	return claims, nil
}
