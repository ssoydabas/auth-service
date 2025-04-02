package service

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AccountServiceTestSuite struct {
	suite.Suite
	mockRepo *MockAccountRepository
	service  service.AccountService
}

func (suite *AccountServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockAccountRepository)
	suite.service = service.NewAccountService(suite.mockRepo)
}

func (suite *AccountServiceTestSuite) TearDownTest() {
	suite.mockRepo.ExpectedCalls = nil
}

func (suite *AccountServiceTestSuite) createTestAccount(id uint, email, phone string) *models.Account {
	return &models.Account{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:              email,
		Phone:              phone,
		VerificationStatus: "pending",
		FirstName:          "John",
		LastName:           "Doe",
	}
}

func (suite *AccountServiceTestSuite) createTestAccountRequest(email, password, phone string) dto.CreateAccountRequest {
	return dto.CreateAccountRequest{
		Email:     email,
		Password:  password,
		Phone:     phone,
		FirstName: "John",
		LastName:  "Doe",
	}
}

func (suite *AccountServiceTestSuite) createTestAuthRequest(email, password, phone string) dto.AuthenticateAccountRequest {
	return dto.AuthenticateAccountRequest{
		Email:    email,
		Password: password,
		Phone:    phone,
	}
}

func (suite *AccountServiceTestSuite) generateTestToken(userID uint) string {
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return token
}

func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
