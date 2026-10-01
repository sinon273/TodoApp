package task_service

import (
	task_ports_repository "TodoApp/internal/features/tasks/ports"
)

type TasksService struct {
	taskRepository task_ports_repository.TaskRepository
}

func NewTasksService(taskRepository task_ports_repository.TaskRepository) *TasksService {
	return &TasksService{taskRepository: taskRepository}
}
