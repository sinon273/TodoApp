package task_cached_repository

import (
	"TodoApp/internal/core/domain"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

func (r *CachedRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	created, err := r.mainRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	model := taskDomainToModel(created)
	bytes, err := json.Marshal(model)
	if err != nil {
		r.log.Error("marshal task for cache", zap.Error(err))
	} else {
		if err = r.pool.Set(ctx, taskKey(created.ID), bytes, r.pool.TTL()).Err(); err != nil {
			r.log.Error("set task to cache", zap.Error(err))
		}
	}

	if err = r.pool.Del(ctx, tasksListKey(nil), tasksListKey(&created.AuthorUserID)).Err(); err != nil {
		r.log.Error("invalidate task lists", zap.Error(err))
	}

	return created, nil
}
