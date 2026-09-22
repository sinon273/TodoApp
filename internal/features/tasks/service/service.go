package task_service

import (
	"TodoApp/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type TasksService struct {
	taskRepository TaskRepository
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *uuid.UUID, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (domain.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
	PatchTask(ctx context.Context, taskID uuid.UUID, taskPatch domain.Task) (domain.Task, error)
}

func NewTasksService(taskRepository TaskRepository) *TasksService {
	return &TasksService{taskRepository: taskRepository}
}
