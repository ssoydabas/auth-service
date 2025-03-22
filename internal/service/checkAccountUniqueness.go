package service

import (
	"context"
	"fmt"

	"github.com/ssoydabas/auth-service/pkg/errors"
)

func (s *accountService) checkAccountUniqueness(ctx context.Context, email, phone string) error {
	emailExists, err := s.accountRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return errors.ConflictError("email already in use")
	}

	phoneExists, err := s.accountRepository.ExistsByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("failed to check phone existence: %w", err)
	}
	if phoneExists {
		return errors.ConflictError("phone number already in use")
	}

	return nil
}
