package service

import (
	"context"
	"fmt"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/repository"

	"github.com/ssoydabas/auth-service/models"
)

type AccountService interface {
	GetAccounts(ctx context.Context, page, pageSize int, search string) (*dto.PaginatedResponse, error)
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) error
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
