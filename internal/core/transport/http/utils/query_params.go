package core_http_utils

import (
	core_error "TodoApp/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	val, err := uuid.Parse(param)
	if err != nil {
		return nil, fmt.Errorf(
			"query param='%s' by key='%s' not a valid uuid: %v: %w",
			param, key, err, core_error.ErrInvalidArgument,
		)
	}
	return &val, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key '%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_error.ErrInvalidArgument,
		)
	}

	return &date, nil
}
