package repository

import (
	"context"
	"time"

	"github.com/ssoydabas/auth-service/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, model models.Account) error
	GetAccountByID(ctx context.Context, id string, preloadTokens bool) (*models.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountByEmailOrPhone(ctx context.Context, email, phone string) (*models.Account, error)
	GetAccountPasswordByAccountID(ctx context.Context, accountID uint) (*models.AccountPassword, error)
	UpdateLastLoginAt(ctx context.Context, accountID uint, lastLoginAt *time.Time) error

	ExistsByEmail(ctx context.Context, email string) bool
	ExistsByPhone(ctx context.Context, phone string) bool

	SetResetPasswordToken(ctx context.Context, accountID uint, token string) error
	GetAccountByResetPasswordToken(ctx context.Context, token string) (*models.Account, error)
	UpdateAccountPassword(ctx context.Context, accountID uint, password string) error

	UpdateAccountVerificationStatus(ctx context.Context, accountID uint, status string) error
	GetAccountByEmailVerificationToken(ctx context.Context, token string) (*models.Account, error)
	ClearEmailVerificationToken(ctx context.Context, accountID uint) error
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

func (r *accountRepository) GetAccountByID(ctx context.Context, id string, preloadTokens bool) (*models.Account, error) {
	var account models.Account
	query := r.db.WithContext(ctx)

	if preloadTokens {
		query = query.Preload("AccountTokens")
	}

	if err := query.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAccountByEmailOrPhone(ctx context.Context, email, phone string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).Where("email = ? OR phone = ?", email, phone).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAccountPasswordByAccountID(ctx context.Context, accountID uint) (*models.AccountPassword, error) {
	var accountPassword models.AccountPassword
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).First(&accountPassword).Error; err != nil {
		return nil, err
	}
	return &accountPassword, nil
}

func (r *accountRepository) ExistsByEmail(ctx context.Context, email string) bool {
	exists := false
	err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Select("1").
		Where("email = ?", email).
		First(&exists).
		Error

	if err == gorm.ErrRecordNotFound {
		return false
	}
	return exists
}

func (r *accountRepository) ExistsByPhone(ctx context.Context, phone string) bool {
	exists := false
	err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Select("1").
		Where("phone = ?", phone).
		First(&exists).
		Error

	if err == gorm.ErrRecordNotFound {
		return false
	}
	return exists
}

func (r *accountRepository) SetResetPasswordToken(ctx context.Context, accountID uint, token string) error {
	return r.db.WithContext(ctx).Model(&models.AccountToken{}).Where("account_id = ?", accountID).Update("reset_password_token", token).Error
}

func (r *accountRepository) GetAccountByResetPasswordToken(ctx context.Context, token string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).
		Preload("AccountTokens").
		Joins("JOIN account_tokens ON account_tokens.account_id = accounts.id").
		Where("account_tokens.reset_password_token = ?", token).
		First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateAccountPassword(ctx context.Context, accountID uint, password string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.AccountPassword{}).
			Where("account_id = ?", accountID).
			Update("password", password).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.AccountToken{}).
			Where("account_id = ?", accountID).
			Update("reset_password_token", "").Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *accountRepository) GetAccountByEmailVerificationToken(ctx context.Context, token string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).
		Preload("AccountTokens").
		Joins("JOIN account_tokens ON account_tokens.account_id = accounts.id").
		Where("account_tokens.email_verification_token = ?", token).
		First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateAccountVerificationStatus(ctx context.Context, accountID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Account{}).
		Where("id = ?", accountID).
		Update("verification_status", status).Error
}

func (r *accountRepository) ClearEmailVerificationToken(ctx context.Context, accountID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.AccountToken{}).
		Where("account_id = ?", accountID).
		Update("email_verification_token", "").Error
}

func (r *accountRepository) UpdateLastLoginAt(ctx context.Context, accountID uint, lastLoginAt *time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Account{}).
		Where("id = ?", accountID).
		Update("last_login_at", lastLoginAt).Error
}
