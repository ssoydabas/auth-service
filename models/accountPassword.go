package models

import (
	"gorm.io/gorm"
)

type AccountPassword struct {
	gorm.Model
	AccountID uint   `json:"account_id" gorm:"unique index"`
	Password  string `json:"password" validate:"required,min=8"`
}
