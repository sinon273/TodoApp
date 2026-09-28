package statistics_transport_http

import (
	"TodoApp/internal/core/domain"
	core_logger "TodoApp/internal/core/logger"
	core_http_response "TodoApp/internal/core/transport/http/response"
	core_http_utils "TodoApp/internal/core/transport/http/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"50"`
	TasksCompleted             int      `json:"tasks_completed" example:"10"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"20"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

// GetStatistics godoc
// @Summary Получение статистки
// @Description Получение статистистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags statistics
// @Produce json
// @Param user_id query uuid false "Фильтрация статистики по конкретному пользователю"
// @Param from query string false "Начало промежутка рассмотрения статистики(включительно), формат: YYYY-MM-DD"
// @Param to query string false "Конец промежутка рассмотрения статистики(не включительно), формат: YYYY-MM-DD"
// @Success 200 {object} GetStatisticsResponse "Успешное получение статистики"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/form/to query params",
		)
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}
	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}

}

func getUserIDFromToQueryParams(r *http.Request) (*uuid.UUID, *time.Time, *time.Time, error) {

	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_utils.GetUUIDQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_utils.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_utils.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
