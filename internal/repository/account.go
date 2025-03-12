package repository

import (
	"context"

	"github.com/ssoydabas/auth-service/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	GetAccounts(ctx context.Context, page, pageSize int) ([]models.Account, int64, error)
	CreateAccount(ctx context.Context, model models.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) GetAccounts(ctx context.Context, page, pageSize int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Account{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (r *accountRepository) CreateAccount(ctx context.Context, model models.Account) error {
	return r.db.WithContext(ctx).Create(&model).Error
}
