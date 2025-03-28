package middleware

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := validator.New()
	err := v.RegisterValidation("id", validateObjectID)
	if err != nil {
		return nil, err
	}

	return &Validator{
		validator: v,
	}, nil
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func validateObjectID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("OBJECTID is not valud %+v", err)
	}
	return err == nil
}

type SchemaFunc = func(c *fiber.Ctx) (interface{}, string, error)

func WithValidation(v *Validator, getSchema SchemaFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schema, name, err := getSchema(c)
		if err != nil {
			fmt.Println("ERROR HERE", err.Error())
			return err
		}

		fmt.Printf(">>>>>>, %+v", schema)

		if err := v.Validate(schema); err != nil {
			return utils.ValidatorError(err)
		}
		if name != "" {
			c.Locals(name, schema)
		}

		return c.Next()
	}
}
