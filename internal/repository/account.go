package repository

import (
	"context"

	"github.com/ssoydabas/auth-service/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, model models.Account) error
	GetAccountByID(ctx context.Context, id string) (*models.Account, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, model models.Account) error {
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *accountRepository) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists := false
	err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Select("1").
		Where("email = ?", email).
		First(&exists).
		Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return exists, err
}

func (r *accountRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	exists := false
	err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Select("1").
		Where("phone = ?", phone).
		First(&exists).
		Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return exists, err
}
