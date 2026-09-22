package task_service

import (
	"TodoApp/internal/core/domain"
	"context"
	"fmt"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("task validation failed: %w", err)
	}
	task, err := s.taskRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("task creation failed: %w", err)
	}

	return task, nil
}
