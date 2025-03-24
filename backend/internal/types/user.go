package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

type CreateUserParams struct {
	FirstName string `validate:"required,alpha,min=2,max=48" json:"firstName"`
	LastName  string `validate:"required,alpha,min=2,max=48" json:"lastName"`
	Email     string `validate:"required,email" json:"email"`
	Password  string `validate:"required,min=7,max=256" json:"password"`
}

func NewUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

type UpdateUserParams struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func (p *UpdateUserParams) ToBsonMap() bson.M {
	values := bson.M{}

	if len(p.FirstName) > 0 {
		values["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		values["lastName"] = p.LastName
	}

	return values
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (param *AuthParams) IsValidPassword(encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(param.Password)) == nil
}
