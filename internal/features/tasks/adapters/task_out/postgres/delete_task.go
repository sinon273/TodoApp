package tasks_postgres_repository

import (
	core_error "TodoApp/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *TasksRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"task with id='%s': %w",
			taskID,
			core_error.ErrNotFound,
		)
	}

	return nil
}
