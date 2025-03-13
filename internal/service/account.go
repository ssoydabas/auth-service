package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ssoydabas/auth-service/models"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) error
	AuthenticateAccount(ctx context.Context, req dto.AuthenticateAccountRequest) (string, error)
	GetAccountByID(ctx context.Context, id string) (*dto.AccountResponse, error)
	GetAccountByToken(ctx context.Context, token string) (*dto.AccountResponse, error)
	SetResetPasswordToken(ctx context.Context, req dto.SetResetPasswordTokenRequest) (string, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	VerifyAccountEmail(ctx context.Context, req dto.VerifyAccountRequest) error
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
	if err := s.checkAccountUniqueness(ctx, req.Email, req.Phone); err != nil {
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

func (s *accountService) AuthenticateAccount(ctx context.Context, req dto.AuthenticateAccountRequest) (string, error) {
	account, err := s.accountRepository.GetAccountByEmailOrPhone(ctx, req.Email, req.Phone)
	if err != nil {
		return "", fmt.Errorf("failed to get account by email or phone: %w", err)
	}

	accountPassword, err := s.accountRepository.GetAccountPasswordByAccountID(ctx, account.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get account password by account id: %w", err)
	}

	if !verifyPassword(req.Password, accountPassword.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": account.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
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

func (s *accountService) GetAccountByToken(ctx context.Context, tokenString string) (*dto.AccountResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	userID := fmt.Sprintf("%.0f", claims["sub"].(float64))

	account, err := s.accountRepository.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &dto.AccountResponse{
		ID:                 account.ID,
		FirstName:          account.FirstName,
		LastName:           account.LastName,
		Email:              account.Email,
		Phone:              account.Phone,
		PhotoUrl:           account.PhotoUrl,
		CreatedAt:          account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          account.UpdatedAt.Format(time.RFC3339),
		VerificationStatus: account.VerificationStatus,
	}, nil
}

func (s *accountService) SetResetPasswordToken(ctx context.Context, req dto.SetResetPasswordTokenRequest) (string, error) {
	account, err := s.accountRepository.GetAccountByEmailOrPhone(ctx, req.Email, req.Phone)
	if err != nil {
		return "", fmt.Errorf("failed to get account by email or phone: %w", err)
	}

	token := uuid.New().String()

	if err := s.accountRepository.SetResetPasswordToken(ctx, account.ID, token); err != nil {
		return "", fmt.Errorf("failed to set update password token: %w", err)
	}

	return token, nil
}

func (s *accountService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	account, err := s.accountRepository.GetAccountByResetPasswordToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("failed to get account by reset password token: %w", err)
	}

	if err := s.accountRepository.UpdateAccountPassword(ctx, account.ID, hashPassword(req.Password)); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

func (s *accountService) VerifyAccountEmail(ctx context.Context, req dto.VerifyAccountRequest) error {
	account, err := s.accountRepository.GetAccountByEmailVerificationToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("failed to get account by email verification token: %w", err)
	}

	if err := s.accountRepository.UpdateAccountVerificationStatus(ctx, account.ID, "verified"); err != nil {
		return fmt.Errorf("failed to update account verification status: %w", err)
	}

	if err := s.accountRepository.ClearEmailVerificationToken(ctx, account.ID); err != nil {
		return fmt.Errorf("failed to clear email verification token: %w", err)
	}

	return nil
}
