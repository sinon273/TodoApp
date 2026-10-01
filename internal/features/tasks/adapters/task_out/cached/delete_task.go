package task_cached_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) Delete(
	ctx context.Context,
	taskID uuid.UUID,
) error {
	var authorUserID uuid.UUID

	bytes, err := r.pool.Get(ctx, taskKey(taskID)).Bytes()
	if err == nil {
		var model TaskModel
		if err = json.Unmarshal(bytes, &model); err == nil {
			authorUserID = model.AuthorUserID
		}
	}

	if authorUserID == uuid.Nil {
		task, err := r.mainRepository.GetTask(ctx, taskID)
		if err != nil {
			return fmt.Errorf("get task by id %v: %w", taskID, err)
		}
		authorUserID = task.AuthorUserID
	}

	if err = r.mainRepository.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("delete task by id %v: %w", taskID, err)
	}

	if err = r.pool.Del(ctx, taskKey(taskID), tasksListKey(nil), tasksListKey(&authorUserID)).Err(); err != nil {
		r.log.Error("delete task by id %v", zap.Error(err))
	}
	return nil
}
