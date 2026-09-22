package users_postgres_repository

import (
	"TodoApp/internal/core/domain"
	core_error "TodoApp/internal/core/errors"
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *UsersRepository) GetUser(ctx context.Context,
	id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number 
	FROM todoapp.users
	WHERE id = $1
`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%s': %w", id, core_error.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil
}
