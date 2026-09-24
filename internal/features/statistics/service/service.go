package statistics_service

import (
	"TodoApp/internal/core/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{statisticsRepository: statisticsRepository}
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}
