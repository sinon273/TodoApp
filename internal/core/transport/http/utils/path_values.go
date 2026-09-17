package core_http_utils

import (
	core_error "TodoApp/internal/core/errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf("no key='%s' in path values: %w", key, core_error.ErrInvalidArgument)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("path value='%s' by key='%s' not a valid uuid: %v: %w",
			pathValue,
			key,
			err,
			core_error.ErrInvalidArgument,
		)
	}

	return val, nil
}
