package models

import (
	"gorm.io/gorm"
)

type AccountToken struct {
	gorm.Model
	AccountID              uint   `json:"account_id"`
	ResetPasswordToken     string `json:"reset_password_token" gorm:"unique index"`
	ResetEmailToken        string `json:"reset_email_token" gorm:"unique index"`
	EmailVerificationToken string `json:"email_verification_token" gorm:"unique index"`
	PhoneVerificationToken string `json:"phone_verification_token" gorm:"unique index"`
}
