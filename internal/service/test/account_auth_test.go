package service

import (
	"context"
	"strings"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/models"
	"github.com/ssoydabas/auth-service/pkg/errors"
	"github.com/stretchr/testify/mock"
)

// TestAuthenticateAccount tests the account authentication functionality
func (suite *AccountServiceTestSuite) TestAuthenticateAccount() {
	tests := []struct {
		name          string
		req           dto.AuthenticateAccountRequest
		setupMocks    func()
		wantToken     bool
		wantErr       bool
		expectedError error
	}{
		{
			name: "account not found",
			req:  suite.createTestAuthRequest("nonexistent@example.com", "password123", "+1234567890"),
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByEmailOrPhone", mock.Anything, "nonexistent@example.com", "+1234567890").
					Return(nil, errors.NotFoundError("Account not found"))
			},
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name: "invalid credentials",
			req:  suite.createTestAuthRequest("test@example.com", "wrongpassword", "+1234567890"),
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				suite.mockRepo.On("GetAccountByEmailOrPhone", mock.Anything, "test@example.com", "+1234567890").
					Return(mockAccount, nil)
				suite.mockRepo.On("GetAccountPasswordByAccountID", mock.Anything, uint(1)).
					Return(&models.AccountPassword{
						Password: service.HashPassword("correctpassword"),
					}, nil)
			},
			wantErr:       true,
			expectedError: errors.AuthError("Invalid credentials"),
		},
		{
			name: "successful authentication",
			req:  suite.createTestAuthRequest("test@example.com", "correctpassword", "+1234567890"),
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				suite.mockRepo.On("GetAccountByEmailOrPhone", mock.Anything, "test@example.com", "+1234567890").
					Return(mockAccount, nil)
				suite.mockRepo.On("GetAccountPasswordByAccountID", mock.Anything, uint(1)).
					Return(&models.AccountPassword{
						Password: service.HashPassword("correctpassword"),
					}, nil)
			},
			wantToken: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMocks()

			token, err := suite.service.AuthenticateAccount(context.Background(), tt.req)

			if tt.wantErr {
				suite.Error(err)
				if tt.expectedError != nil {
					suite.Equal(tt.expectedError.Error(), err.Error())
				}
			} else {
				suite.NoError(err)
				if tt.wantToken {
					suite.NotEmpty(token)
					// Verify token format
					suite.Contains(token, ".")
					// Basic JWT format check (header.payload.signature)
					parts := strings.Split(token, ".")
					suite.Len(parts, 3)
				}
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// TestGetAccountByToken tests the account retrieval by token functionality
func (suite *AccountServiceTestSuite) TestGetAccountByToken() {
	tests := []struct {
		name          string
		token         string
		setupMocks    func()
		wantAccount   bool
		wantErr       bool
		expectedError error
	}{
		{
			name:  "invalid token format",
			token: "invalid-token",
			setupMocks: func() {
				// No mocks needed for invalid token format
			},
			wantErr:       true,
			expectedError: errors.BadRequestError("Invalid token"),
		},
		{
			name:  "expired token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImV4cCI6MTcwOTc2ODAwMH0.INVALID_SIGNATURE",
			setupMocks: func() {
				// No mocks needed for expired token
			},
			wantErr:       true,
			expectedError: errors.BadRequestError("Invalid token"),
		},
		{
			name:  "account not found",
			token: suite.generateTestToken(999),
			setupMocks: func() {
				suite.mockRepo.On("GetAccountByID", mock.Anything, "999", false).
					Return(nil, errors.NotFoundError("Account not found"))
			},
			wantErr:       true,
			expectedError: errors.NotFoundError("Account not found"),
		},
		{
			name:  "successful account retrieval",
			token: suite.generateTestToken(1),
			setupMocks: func() {
				mockAccount := suite.createTestAccount(1, "test@example.com", "+1234567890")
				mockAccount.VerificationStatus = "verified"
				suite.mockRepo.On("GetAccountByID", mock.Anything, "1", false).
					Return(mockAccount, nil)
			},
			wantAccount: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.setupMocks()

			account, err := suite.service.GetAccountByToken(context.Background(), tt.token)

			if tt.wantErr {
				suite.Error(err)
				if tt.expectedError != nil {
					suite.Equal(tt.expectedError.Error(), err.Error())
				}
			} else {
				suite.NoError(err)
				if tt.wantAccount {
					suite.NotNil(account)
					suite.Equal(uint(1), account.ID)
					suite.Equal("test@example.com", account.Email)
					suite.Equal("verified", account.VerificationStatus)
				}
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}
