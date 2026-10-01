package task_cached_repository

import (
	"TodoApp/internal/core/domain"
	core_redis_pool "TodoApp/internal/core/repository/redis/pool"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) GetTask(ctx context.Context, taskID uuid.UUID) (domain.Task, error) {
	bytes, err := r.pool.Get(ctx, taskKey(taskID)).Bytes()
	if err == nil {
		var model TaskModel
		if err = json.Unmarshal(bytes, &model); err != nil {
			r.log.Error("unmarshal task from cache", zap.Error(err))
		} else {
			return taskModelToDomain(model), nil
		}
	} else if !errors.Is(err, core_redis_pool.NotFound) {
		r.log.Error("get task from redis", zap.Error(err))
	}

	task, err := r.mainRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, err
	}
	model := taskDomainToModel(task)
	bytesToCache, err := json.Marshal(model)
	if err != nil {
		r.log.Error("marshal task for cache", zap.Error(err))
	} else {
		if err = r.pool.Set(ctx, taskKey(task.ID), bytesToCache, r.pool.TTL()).Err(); err != nil {
			r.log.Error("set task to cache", zap.Error(err))
		}
	}
	return task, nil
}
