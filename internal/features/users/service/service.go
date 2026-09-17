package users_service

import (
	"TodoApp/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type UsersService struct {
	usersRepository UsersRepository
}

func NewUserService(usersRepository UsersRepository) *UsersService {
	return &UsersService{usersRepository: usersRepository}
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user domain.User,
	) (domain.User, error)
}
