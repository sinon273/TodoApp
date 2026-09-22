package task_service

import (
	"TodoApp/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) PatchTask(ctx context.Context, taskID uuid.UUID, taskPatch domain.TaskPatch) (domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(taskPatch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}
	patchedTask, err := s.taskRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
