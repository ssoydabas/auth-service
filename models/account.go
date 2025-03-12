package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	FirstName          string `json:"first_name" validate:"required,min=2,max=50"`
	LastName           string `json:"last_name" validate:"required,min=2,max=50"`
	Email              string `json:"email" validate:"required,email" gorm:"unique index"`
	Phone              string `json:"phone" validate:"required,e164" gorm:"unique index"`
	PhotoUrl           string `json:"photo_url" validate:"omitempty,url"`
	VerificationStatus string `json:"verification_status" validate:"required,oneof=pending verified"`

	AccountPassword AccountPassword `json:"account_password" gorm:"foreignKey:AccountID"`
	AccountTokens   AccountToken    `json:"account_tokens" gorm:"foreignKey:AccountID"`
}
