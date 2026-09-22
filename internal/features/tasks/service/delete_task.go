package task_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) Delete(ctx context.Context, taskID uuid.UUID) error {
	if err := s.taskRepository.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}
	return nil
}
