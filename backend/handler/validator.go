package handler

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	err := v.RegisterValidation("objectId", validateObjectID)
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
