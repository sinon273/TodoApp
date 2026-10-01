package tasks_postgres_repository

import (
	"TodoApp/internal/core/domain"
	core_error "TodoApp/internal/core/errors"
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks(id, title, description, completed, created_at,completed_at, author_user_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserId,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user with id='%d': %w ",
				err,
				task.AuthorUserID,
				core_error.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
