package task_service

import (
	"TodoApp/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) GetTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task fron repository: %w", err)
	}

	return task, nil
}
