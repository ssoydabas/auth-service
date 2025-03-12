package service

import (
	"context"
	"fmt"
)

func (s *accountService) checkAccountUniqueness(ctx context.Context, email, phone string) error {
	emailExists, err := s.accountRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return fmt.Errorf("email already in use")
	}

	phoneExists, err := s.accountRepository.ExistsByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("failed to check phone existence: %w", err)
	}
	if phoneExists {
		return fmt.Errorf("phone number already in use")
	}

	return nil
}
