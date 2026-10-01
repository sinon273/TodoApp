package task_cached_repository

import (
	"TodoApp/internal/core/domain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) PatchTask(
	ctx context.Context,
	taskID uuid.UUID,
	taskPatch domain.Task,
) (domain.Task, error) {
	updated, err := r.mainRepository.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}
	model := taskDomainToModel(updated)
	bytes, err := json.Marshal(model)
	if err != nil {
		r.log.Error("marshal task for cache", zap.Error(err))
	} else {
		if err = r.pool.Set(ctx, taskKey(updated.ID), bytes, r.pool.TTL()).Err(); err != nil {
			r.log.Error("set task to cache", zap.Error(err))
		}
	}
	if err = r.pool.Del(ctx, tasksListKey(nil), tasksListKey(&updated.AuthorUserID)).Err(); err != nil {
		r.log.Error("invalidate task lists", zap.Error(err))
	}

	return updated, nil
}
