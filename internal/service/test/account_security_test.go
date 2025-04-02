package service

import (
	"context"
	"fmt"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/models"
	"github.com/ssoydabas/auth-service/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func (suite *AccountServiceTestSuite) TestPasswordResetFlow() {
	tests := []struct {
		name          string
		setupMocks    func()
		req           dto.SetResetPasswordTokenRequest
		resetReq      dto.ResetPasswordRequest
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful password reset flow",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				suite.mockRepo.On("GetAccountByEmailOrPhone", mock.Anything, "test@example.com", "+1234567890").
					Return(mockAccount, nil)
				suite.mockRepo.On("SetResetPasswordToken", mock.Anything, uint(1), mock.Anything).
					Return(nil)

				suite.mockRepo.On("GetAccountByResetPasswordToken", mock.Anything, mock.Anything).
					Return(mockAccount, nil)
				suite.mockRepo.On("UpdateAccountPassword", mock.Anything, uint(1), mock.Anything).
					Return(nil)
			},
			req: dto.SetResetPasswordTokenRequest{
				Email: "test@example.com",
				Phone: "+1234567890",
			},
			resetReq: dto.ResetPasswordRequest{
				Token:    "valid-token",
				Password: "newSecurePassword123!",
			},
			wantErr: false,
		},
		{
			name: "account not found for reset token request",
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByEmailOrPhone", mock.Anything, "nonexistent@example.com", "+1234567890").
					Return(nil, errors.NotFoundError("Account not found"))
			},
			req: dto.SetResetPasswordTokenRequest{
				Email: "nonexistent@example.com",
				Phone: "+1234567890",
			},
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name: "invalid reset token",
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByResetPasswordToken", mock.Anything, "invalid-token").
					Return(nil, errors.NotFoundError("Account not found"))
			},
			resetReq: dto.ResetPasswordRequest{
				Token:    "invalid-token",
				Password: "newSecurePassword123!",
			},
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name: "repository error during password update",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				suite.mockRepo.On("GetAccountByResetPasswordToken", mock.Anything, "valid-token").
					Return(mockAccount, nil)
				suite.mockRepo.On("UpdateAccountPassword", mock.Anything, uint(1), mock.Anything).
					Return(errors.InternalError(fmt.Errorf("database error")))
			},
			resetReq: dto.ResetPasswordRequest{
				Token:    "valid-token",
				Password: "newSecurePassword123!",
			},
			wantErr:       true,
			expectedError: errors.InternalError(fmt.Errorf("database error")),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			if tt.req.Email != "" {
				token, err := suite.service.SetResetPasswordToken(context.Background(), tt.req)
				if tt.wantErr {
					suite.Error(err)
					if tt.expectedError != nil {
						suite.Equal(tt.expectedError.Error(), err.Error())
					}
					return
				}
				suite.NoError(err)
				suite.NotEmpty(token)
			}

			if tt.resetReq.Token != "" {
				err := suite.service.ResetPassword(context.Background(), tt.resetReq)
				if tt.wantErr {
					suite.Error(err)
					if tt.expectedError != nil {
						suite.Equal(tt.expectedError.Error(), err.Error())
					}
				} else {
					suite.NoError(err)
				}
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *AccountServiceTestSuite) TestEmailVerificationFlow() {
	tests := []struct {
		name          string
		setupMocks    func()
		req           dto.VerifyAccountRequest
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful email verification",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "pending"
				suite.mockRepo.On("GetAccountByEmailVerificationToken", mock.Anything, "valid-token").
					Return(mockAccount, nil)
				suite.mockRepo.On("UpdateAccountVerificationStatus", mock.Anything, uint(1), "verified").
					Return(nil)
				suite.mockRepo.On("ClearEmailVerificationToken", mock.Anything, uint(1)).
					Return(nil)
			},
			req: dto.VerifyAccountRequest{
				Token: "valid-token",
			},
			wantErr: false,
		},
		{
			name: "invalid verification token",
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByEmailVerificationToken", mock.Anything, "invalid-token").
					Return(nil, errors.NotFoundError("Account not found"))
			},
			req: dto.VerifyAccountRequest{
				Token: "invalid-token",
			},
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name: "already verified account",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "verified"
				suite.mockRepo.On("GetAccountByEmailVerificationToken", mock.Anything, "valid-token").
					Return(mockAccount, nil)
				suite.mockRepo.On("UpdateAccountVerificationStatus", mock.Anything, uint(1), "verified").
					Return(errors.BadRequestError("Account already verified"))
			},
			req: dto.VerifyAccountRequest{
				Token: "valid-token",
			},
			wantErr:       true,
			expectedError: errors.BadRequestError("Account already verified"),
		},
		{
			name: "repository error during verification",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "pending"
				suite.mockRepo.On("GetAccountByEmailVerificationToken", mock.Anything, "valid-token").
					Return(mockAccount, nil)
				suite.mockRepo.On("UpdateAccountVerificationStatus", mock.Anything, uint(1), "verified").
					Return(errors.InternalError(fmt.Errorf("database error")))
			},
			req: dto.VerifyAccountRequest{
				Token: "valid-token",
			},
			wantErr:       true,
			expectedError: errors.InternalError(fmt.Errorf("database error")),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			err := suite.service.VerifyAccountEmail(context.Background(), tt.req)
			if tt.wantErr {
				suite.Error(err)
				if tt.expectedError != nil {
					suite.Equal(tt.expectedError.Error(), err.Error())
				}
			} else {
				suite.NoError(err)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *AccountServiceTestSuite) TestGetEmailVerificationToken() {
	tests := []struct {
		name          string
		setupMocks    func()
		accountID     string
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful token retrieval",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "pending"
				mockAccount.AccountTokens = models.AccountToken{
					EmailVerificationToken: "test-verification-token",
				}
				suite.mockRepo.On("GetAccountByID", mock.Anything, "1", true).
					Return(mockAccount, nil)
			},
			accountID: "1",
			wantErr:   false,
		},
		{
			name: "account not found",
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByID", mock.Anything, "999", true).
					Return(nil, errors.NotFoundError("Account not found"))
			},
			accountID:     "999",
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name: "already verified account",
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "verified"
				suite.mockRepo.On("GetAccountByID", mock.Anything, "1", true).
					Return(mockAccount, nil)
			},
			accountID:     "1",
			wantErr:       true,
			expectedError: errors.BadRequestError("Account already verified"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			token, err := suite.service.GetAccountEmailVerificationTokenByID(context.Background(), tt.accountID)
			if tt.wantErr {
				suite.Error(err)
				if tt.expectedError != nil {
					suite.Equal(tt.expectedError.Error(), err.Error())
				}
			} else {
				suite.NoError(err)
				suite.NotEmpty(token)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}
