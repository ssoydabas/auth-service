package dto

import (
	"github.com/ssoydabas/auth-service/pkg/validator"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (r *CreateAccountRequest) Validate() error {
	return validator.ValidateStruct(r)
}
