package tasks_http

import (
	"TodoApp/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type TasksDTOResponse struct {
	ID           uuid.UUID  `json:"id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	Version      int        `json:"version" example:"2"`
	Title        string     `json:"title" example:"Домашнее задание"`
	Description  *string    `json:"description" example:"Сделать до понедельника домашнее задание"`
	Completed    bool       `json:"completed" example:"false"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-02-26T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"null"`
	AuthorUserID uuid.UUID  `json:"author_user_id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d468"`
}

func taskDTOFromDomain(task domain.Task) TasksDTOResponse {
	return TasksDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TasksDTOResponse {
	dtos := make([]TasksDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}
	return dtos
}
