package dto

import (
	"fmt"

	"github.com/ssoydabas/auth-service/pkg/validator"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,min=8"`
}

type AuthenticateAccountRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
	Password string `json:"password" validate:"required,min=8"`
}

type SetResetPasswordTokenRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone" validate:"omitempty,e164"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

func (r *CreateAccountRequest) Validate() error {
	return validator.ValidateStruct(r)
}

func (r *AuthenticateAccountRequest) Validate() error {
	if r.Email == "" && r.Phone == "" {
		return fmt.Errorf("either email or phone must be provided")
	}

	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	return validator.ValidateStruct(r)
}

func (r *SetResetPasswordTokenRequest) Validate() error {
	if r.Email == "" && r.Phone == "" {
		return fmt.Errorf("either email or phone must be provided")
	}

	return validator.ValidateStruct(r)
}

func (r *ResetPasswordRequest) Validate() error {
	if r.Token == "" {
		return fmt.Errorf("token is required")
	}

	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	if r.Password != r.ConfirmPassword {
		return fmt.Errorf("password and confirm password do not match")
	}

	return validator.ValidateStruct(r)
}
