package tasks_postgres_repository

import (
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	"context"
)

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func (r *TasksRepository) rGetTasks(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func NewTaskRepository(pool core_postgres_pool.Pool) *TasksRepository {
	return &TasksRepository{pool: pool}
}
