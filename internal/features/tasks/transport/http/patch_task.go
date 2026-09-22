package tasks_transport

import (
	"TodoApp/internal/core/domain"
	core_logger "TodoApp/internal/core/logger"
	core_http_request "TodoApp/internal/core/transport/http/request"
	core_http_response "TodoApp/internal/core/transport/http/response"
	core_http_types "TodoApp/internal/core/transport/http/types"
	core_http_utils "TodoApp/internal/core/transport/http/utils"
	"fmt"
	"net/http"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("'Title' length must be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("'Description' length must be between 1 and 100")
			}
		}
	}
	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("'Completed' can't be NULL")
		}
	}
	return nil
}

type PatchUserResponse TasksDTOResponse

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get path value")
		return
	}

	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
