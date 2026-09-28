package users_transport_http

import (
	core_logger "TodoApp/internal/core/logger"
	core_http_response "TodoApp/internal/core/transport/http/response"
	core_http_utils "TodoApp/internal/core/transport/http/utils"
	"net/http"
)

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление существующего в системе пользователя по его ID
// @Tags users
// @Param id path uuid true "ID удаляемого пользователя"
// @Success 204 "Успешное удаление пользователя"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}
	err = h.usersService.DeleteUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}
	responseHandler.NoContentResponse()
}
