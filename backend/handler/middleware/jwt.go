package middleware

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	fmt.Println("TOKEN FIRST ITEM", token[0])
	claims, err := validateToken(token[0])
	if err != nil {
		return fmt.Errorf("token is invalid")
	}

	ctx := context.WithValue(c.Context(), "userID", claims["id"])
	c.SetUserContext(ctx)

	fmt.Println("Token: ", claims)

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", t.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorize")
	}

	return claims, nil
}
