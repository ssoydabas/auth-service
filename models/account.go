package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	FirstName          string     `json:"first_name" validate:"required,min=2,max=50"`
	LastName           string     `json:"last_name" validate:"required,min=2,max=50"`
	Email              string     `json:"email" validate:"required,email" gorm:"uniqueIndex:idx_email"`
	Phone              string     `json:"phone" validate:"required,e164" gorm:"uniqueIndex:idx_phone"`
	PhotoUrl           string     `json:"photo_url" validate:"omitempty,url"`
	VerificationStatus string     `json:"verification_status" validate:"required,oneof=pending verified"`
	Role               string     `json:"role" validate:"required,oneof=common admin manager teacher student" gorm:"default:common"`
	LastLoginAt        *time.Time `json:"last_login_at"`

	AccountPassword AccountPassword `json:"account_password" gorm:"foreignKey:AccountID"`
	AccountTokens   AccountToken    `json:"account_tokens" gorm:"foreignKey:AccountID"`
}
