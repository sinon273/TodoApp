package tasks_transport

import (
	"TodoApp/internal/core/domain"
	core_logger "TodoApp/internal/core/logger"
	core_http_request "TodoApp/internal/core/transport/http/request"
	core_http_response "TodoApp/internal/core/transport/http/response"
	"net/http"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title        string    `json:"title" validate:"required,min=1,max=100"`
	Description  *string   `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID uuid.UUID `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TasksDTOResponse

func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
