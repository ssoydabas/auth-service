package service

import (
	"context"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/models"
	"github.com/stretchr/testify/mock"
)

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, model models.Account) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockAccountRepository) AuthenticateAccount(ctx context.Context, model dto.AuthenticateAccountRequest) (string, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByID(ctx context.Context, id string, preloadTokens bool) (*models.Account, error) {
	args := m.Called(ctx, id, preloadTokens)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByEmailOrPhone(ctx context.Context, email, phone string) (*models.Account, error) {
	args := m.Called(ctx, email, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountPasswordByAccountID(ctx context.Context, accountID uint) (*models.AccountPassword, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AccountPassword), args.Error(1)
}

func (m *MockAccountRepository) ExistsByEmail(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0)
}

func (m *MockAccountRepository) ExistsByPhone(ctx context.Context, phone string) bool {
	args := m.Called(ctx, phone)
	return args.Bool(0)
}

func (m *MockAccountRepository) SetResetPasswordToken(ctx context.Context, accountID uint, token string) error {
	args := m.Called(ctx, accountID, token)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountByResetPasswordToken(ctx context.Context, token string) (*models.Account, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateAccountPassword(ctx context.Context, accountID uint, password string) error {
	args := m.Called(ctx, accountID, password)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateAccountVerificationStatus(ctx context.Context, accountID uint, status string) error {
	args := m.Called(ctx, accountID, status)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountByEmailVerificationToken(ctx context.Context, token string) (*models.Account, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) ClearEmailVerificationToken(ctx context.Context, accountID uint) error {
	args := m.Called(ctx, accountID)
	return args.Error(0)
}
