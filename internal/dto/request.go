package dto

import (
	"github.com/go-playground/validator/v10"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (r *CreateAccountRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
