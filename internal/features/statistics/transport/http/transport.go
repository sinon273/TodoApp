package statistics_transport_http

import (
	"TodoApp/internal/core/domain"
	core_http_server "TodoApp/internal/core/transport/http/server"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

func NewStatisticsHTTPHandler(statisticsService StatisticsService) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{statisticsService: statisticsService}
}

type StatisticsService interface {
	GetStatistics(ctx context.Context, userID *uuid.UUID, from *time.Time, to *time.Time) (domain.Statistics, error)
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
