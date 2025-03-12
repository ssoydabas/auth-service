package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/repository"

	"github.com/ssoydabas/auth-service/models"

	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) error
	GetAccountByID(ctx context.Context, id string) (*dto.AccountResponse, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) error {
	if err := s.checkUniqueness(ctx, req.Email, req.Phone); err != nil {
		return err
	}

	account := models.Account{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		Phone:              req.Phone,
		VerificationStatus: "pending",
		AccountPassword: models.AccountPassword{
			Password: hashPassword(req.Password),
		},
		AccountTokens: models.AccountToken{
			EmailVerificationToken: uuid.New().String(),
			PhoneVerificationToken: uuid.New().String(),
		},
	}

	if err := s.accountRepository.CreateAccount(ctx, account); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (s *accountService) GetAccountByID(ctx context.Context, id string) (*dto.AccountResponse, error) {
	account, err := s.accountRepository.GetAccountByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by id: %w", err)
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
		Phone:     account.Phone,
		PhotoUrl:  account.PhotoUrl,
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
		UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
	}

	return &response, nil
}
