package task_cached_repository

import (
	"strconv"

	"github.com/google/uuid"
)

func taskKey(id uuid.UUID) string {
	return "task:" + id.String()
}

func tasksListKey(userID *uuid.UUID) string {
	if userID == nil {
		return "tasks:all"
	}

	return "tasks:" + userID.String()
}

func tasksListField(limit, offset *int) string {
	limitStr := "nil"
	if limit != nil {
		limitStr = strconv.Itoa(*limit)
	}

	offsetStr := "nil"
	if offset != nil {
		offsetStr = strconv.Itoa(*offset)
	}

	return limitStr + ":" + offsetStr
}
