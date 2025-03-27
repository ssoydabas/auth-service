package handler

import (
	"github.com/ssoydabas/auth-service/internal/dto"
)

func mapAccountToResponse(account dto.AccountResponse) dto.AccountResponse {
	return dto.AccountResponse{
		ID:                 account.ID,
		FirstName:          account.FirstName,
		LastName:           account.LastName,
		Email:              account.Email,
		Phone:              account.Phone,
		PhotoUrl:           account.PhotoUrl,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
		VerificationStatus: account.VerificationStatus,
	}
}
