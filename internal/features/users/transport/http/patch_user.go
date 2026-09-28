package users_transport_http

import (
	"TodoApp/internal/core/domain"
	core_logger "TodoApp/internal/core/logger"
	core_http_request "TodoApp/internal/core/transport/http/request"
	core_http_response "TodoApp/internal/core/transport/http/response"
	core_http_types "TodoApp/internal/core/transport/http/types"
	core_http_utils "TodoApp/internal/core/transport/http/utils"
	"fmt"
	"net/http"
	"strings"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Денис Денисоввич"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+78005553535"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` length must be between 3 and 100")
		}
	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`PhoneNumber` length must be between 10 and 15")
		}

		if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
			return fmt.Errorf("`PhoneNumber` must startswith '+'")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в бд не меняется
// @Description 2. **Явно передано значение**: `"phone_number": "+78005553535"` - устанавливает новый номер телефона в бд
// @Description 3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description 4. Ограничение: `full_name` не может быть выставлен как null
// @Tags users
// @Accept json
// @Produce json
// @Param id path uuid true "ID изменяемого пользователя"
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Успешно изменённый пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}
	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf("PatchUserRequest fields:\nFullName: '%v'\nPhoneNumber: '%v'",
			request.FullName,
			request.PhoneNumber,
		),
	)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
