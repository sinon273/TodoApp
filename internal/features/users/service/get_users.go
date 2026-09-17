package users_service

import (
	"TodoApp/internal/core/domain"
	core_error "TodoApp/internal/core/errors"
	"context"
	"fmt"
)

func (s *UsersService) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be a non-negative integer: %w", core_error.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be a non-negative integer: %w", core_error.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}
