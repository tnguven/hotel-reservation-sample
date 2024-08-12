package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	err := v.RegisterValidation("id", validateObjectID)
	if err != nil {
		panic(err)
	}

	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func validateObjectID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

type SchemaFunc = func(c *fiber.Ctx) (interface{}, error)

func MiddlewareValidation(v *Validator, getSchema SchemaFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schema, err := getSchema(c)

		if err != nil {
			return err
		}
		if err := v.Validate(schema); err != nil {
			return err

		}

		return c.Next()
	}
}
