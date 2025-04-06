package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/models"
	"github.com/ssoydabas/auth-service/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func (suite *AccountServiceTestSuite) TestCreateAccount() {
	tests := []struct {
		name          string
		req           dto.CreateAccountRequest
		setupMocks    func()
		wantToken     bool
		wantErr       bool
		expectedError error
	}{
		{
			name: "email already exists",
			req:  suite.createTestAccountRequest("test@example.com", "password123", "+1234567890"),
			setupMocks: func() {
				suite.mockRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(true)
			},
			wantErr:       true,
			expectedError: errors.ConflictError("email already in use"),
		},
		{
			name: "phone already exists",
			req:  suite.createTestAccountRequest("test@example.com", "password123", "+1234567890"),
			setupMocks: func() {
				suite.mockRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false)
				suite.mockRepo.On("ExistsByPhone", mock.Anything, "+1234567890").Return(true)
			},
			wantErr:       true,
			expectedError: errors.ConflictError("phone number already in use"),
		},
		{
			name: "repository error during creation",
			req:  suite.createTestAccountRequest("test@example.com", "password123", "+1234567890"),
			setupMocks: func() {
				suite.mockRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false)
				suite.mockRepo.On("ExistsByPhone", mock.Anything, "+1234567890").Return(false)
				suite.mockRepo.On("CreateAccount", mock.Anything, mock.Anything).Return(fmt.Errorf("db error"))
			},
			wantErr:       true,
			expectedError: fmt.Errorf("db error"),
		},
		{
			name: "successful account creation",
			req:  suite.createTestAccountRequest("test@example.com", "password123", "+1234567890"),
			setupMocks: func() {
				suite.mockRepo.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false)
				suite.mockRepo.On("ExistsByPhone", mock.Anything, "+1234567890").Return(false)
				suite.mockRepo.On("CreateAccount", mock.Anything, mock.MatchedBy(func(acc models.Account) bool {
					return acc.Email == "test@example.com" &&
						acc.Phone == "+1234567890" &&
						acc.VerificationStatus == "pending" &&
						acc.AccountTokens.EmailVerificationToken != "" &&
						acc.AccountTokens.PhoneVerificationToken != ""
				})).Return(nil)
			},
			wantToken: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			token, err := suite.service.CreateAccount(context.Background(), tt.req)

			if tt.wantErr {
				suite.Error(err)
				if tt.expectedError != nil {
					suite.Equal(tt.expectedError.Error(), err.Error())
				}
			} else {
				suite.NoError(err)
				if tt.wantToken {
					suite.NotEmpty(token)
					_, err := uuid.Parse(token)
					suite.NoError(err)
				}
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}
