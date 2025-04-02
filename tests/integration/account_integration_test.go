package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/repository"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/pkg/config"
	pkgerrors "github.com/ssoydabas/auth-service/pkg/errors"
	"github.com/ssoydabas/auth-service/pkg/postgres"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find project root (go.mod file)")
		}
		dir = parent
	}
}

type AccountIntegrationTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service service.AccountService
	ctx     context.Context
}

func (suite *AccountIntegrationTestSuite) SetupSuite() {
	projectRoot, err := findProjectRoot()
	suite.Require().NoError(err)
	suite.Require().NotEmpty(projectRoot)

	testEnvPath := filepath.Join(projectRoot, ".env.test")
	os.Setenv("ENV_FILE", testEnvPath)

	_, err = os.Stat(testEnvPath)
	suite.Require().NoError(err, "Test environment file not found at %s", testEnvPath)

	cfg, err := config.LoadConfig()
	suite.Require().NoError(err)
	suite.Require().NotNil(cfg)

	db, err := postgres.ConnectPQ(*cfg)
	suite.Require().NoError(err)
	suite.Require().NotNil(db)
	suite.db = db

	accountRepo := repository.NewAccountRepository(db)
	suite.service = service.NewAccountService(accountRepo)

	suite.ctx = context.Background()
}

func (suite *AccountIntegrationTestSuite) TearDownSuite() {
	sqlDB, err := suite.db.DB()
	suite.Require().NoError(err)
	err = sqlDB.Close()
	suite.Require().NoError(err)
}

func (suite *AccountIntegrationTestSuite) SetupTest() {
	// We truncate all the tables before tests
	err := suite.db.Exec(`
		DO $$
		DECLARE
			table_name text;
		BEGIN
			FOR table_name IN
				SELECT tablename
				FROM pg_tables
				WHERE schemaname = 'public'
			LOOP
				EXECUTE 'TRUNCATE TABLE ' || table_name || ' CASCADE';
			END LOOP;
		END $$;
	`).Error
	suite.Require().NoError(err)
}

func (suite *AccountIntegrationTestSuite) TestCreateAndAuthenticateAccount() {
	createReq := dto.CreateAccountRequest{
		Email:     "test@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	token, err := suite.service.CreateAccount(suite.ctx, createReq)
	suite.NoError(err)
	suite.NotEmpty(token)

	authReq := dto.AuthenticateAccountRequest{
		Email:    "test@example.com",
		Password: "password123",
		Phone:    "+1234567890",
	}

	authToken, err := suite.service.AuthenticateAccount(suite.ctx, authReq)
	suite.NoError(err)
	suite.NotEmpty(authToken)

	authReq.Password = "wrongpassword"
	_, err = suite.service.AuthenticateAccount(suite.ctx, authReq)
	suite.Error(err)
	suite.IsType(pkgerrors.AuthError(""), err)
}

func (suite *AccountIntegrationTestSuite) TestDuplicateEmailAndPhone() {
	createReq := dto.CreateAccountRequest{
		Email:     "test@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	token, err := suite.service.CreateAccount(suite.ctx, createReq)
	suite.NoError(err)
	suite.NotEmpty(token)

	createReq.FirstName = "Jane"
	_, err = suite.service.CreateAccount(suite.ctx, createReq)
	suite.Error(err)
	suite.IsType(pkgerrors.ConflictError(""), err)

	createReq.Email = "different@example.com"
	createReq.Phone = "+1234567890"
	_, err = suite.service.CreateAccount(suite.ctx, createReq)
	suite.Error(err)
	suite.IsType(pkgerrors.ConflictError(""), err)
}

func (suite *AccountIntegrationTestSuite) TestEmailVerificationFlow() {
	createReq := dto.CreateAccountRequest{
		Email:     "test@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	token, err := suite.service.CreateAccount(suite.ctx, createReq)
	suite.NoError(err)
	suite.NotEmpty(token)

	account, err := suite.service.GetAccountByEmail(suite.ctx, createReq.Email)
	suite.NoError(err)
	suite.NotEmpty(account)

	verificationToken, err := suite.service.GetAccountEmailVerificationTokenByID(suite.ctx, fmt.Sprintf("%d", account.ID))
	suite.NoError(err)
	suite.NotEmpty(verificationToken)

	verifyReq := dto.VerifyAccountRequest{
		Token: verificationToken,
	}

	err = suite.service.VerifyAccountEmail(suite.ctx, verifyReq)
	suite.NoError(err)

	err = suite.service.VerifyAccountEmail(suite.ctx, verifyReq)
	suite.Error(err)
	suite.IsType(pkgerrors.NotFoundError(""), err)
}

func (suite *AccountIntegrationTestSuite) TestPasswordResetFlow() {
	createReq := dto.CreateAccountRequest{
		Email:     "test@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	token, err := suite.service.CreateAccount(suite.ctx, createReq)
	suite.NoError(err)
	suite.NotEmpty(token)

	resetTokenReq := dto.SetResetPasswordTokenRequest{
		Email: "test@example.com",
		Phone: "+1234567890",
	}

	resetToken, err := suite.service.SetResetPasswordToken(suite.ctx, resetTokenReq)
	suite.NoError(err)
	suite.NotEmpty(resetToken)

	resetReq := dto.ResetPasswordRequest{
		Token:    resetToken,
		Password: "newpassword123",
	}

	err = suite.service.ResetPassword(suite.ctx, resetReq)
	suite.NoError(err)

	// We try to authenticate with old password
	authReq := dto.AuthenticateAccountRequest{
		Email:    "test@example.com",
		Password: "password123",
		Phone:    "+1234567890",
	}

	_, err = suite.service.AuthenticateAccount(suite.ctx, authReq)
	suite.Error(err)
	suite.IsType(pkgerrors.AuthError(""), err)

	// We authenticate with new password
	authReq.Password = "newpassword123"
	authToken, err := suite.service.AuthenticateAccount(suite.ctx, authReq)
	suite.NoError(err)
	suite.NotEmpty(authToken)

	// Try to use the same reset token again
	err = suite.service.ResetPassword(suite.ctx, resetReq)
	suite.Error(err)
	suite.IsType(pkgerrors.NotFoundError(""), err)
}

func TestAccountIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AccountIntegrationTestSuite))
}
