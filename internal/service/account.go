package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/repository"

	"github.com/ssoydabas/auth-service/models"
)

type AccountService interface {
	GetAccounts(ctx context.Context, page, pageSize int, search string) (*dto.PaginatedResponse, error)
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

func (s *accountService) GetAccounts(ctx context.Context, page, pageSize int, search string) (*dto.PaginatedResponse, error) {
	return &dto.PaginatedResponse{
		Data:        []string{},
		CurrentPage: 1,
		PageSize:    10,
		TotalItems:  0,
		TotalPages:  0,
	}, nil
}

func (s *accountService) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) error {
	if err := s.checkUniqueness(ctx, req.Email, req.Phone); err != nil {
		return err
	}

	model := models.Account{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		Phone:              req.Phone,
		VerificationStatus: "pending",
	}

	if err := s.accountRepository.CreateAccount(ctx, model); err != nil {
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
