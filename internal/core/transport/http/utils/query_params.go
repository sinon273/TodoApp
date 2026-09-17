package core_http_utils

import (
	core_error "TodoApp/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid inteeger: %v: %w",
			param, key, err, core_error.ErrInvalidArgument)
	}

	return &val, nil
}
