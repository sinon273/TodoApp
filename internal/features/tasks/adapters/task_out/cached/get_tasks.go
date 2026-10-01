package task_cached_repository

import (
	"TodoApp/internal/core/domain"
	core_redis_pool "TodoApp/internal/core/repository/redis/pool"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) GetTasks(
	ctx context.Context,
	userID *uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Task, error) {

	key := tasksListKey(userID)
	field := tasksListField(limit, offset)
	bytes, err := r.pool.HGet(ctx, key, field).Bytes()
	if err == nil {
		var models []TaskModel

		if err = json.Unmarshal(bytes, &models); err != nil {
			r.log.Error("unmarshal tasks from cache", zap.Error(err))

		} else {
			return tasksModelsToDomains(models), nil
		}

	} else if !errors.Is(err, core_redis_pool.NotFound) {
		r.log.Error("get tasks from redis", zap.Error(err))
	}

	tasks, err := r.mainRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from main repository: %w", err)
	}

	models := make([]TaskModel, len(tasks))
	for i, task := range tasks {
		models[i] = taskDomainToModel(task)
	}

	bytesToCache, err := json.Marshal(models)
	if err != nil {
		r.log.Error("marshal tasks for cache", zap.Error(err))
	} else {
		if err = r.pool.HSet(ctx, key, field, bytesToCache).Err(); err != nil {
			r.log.Error("set tasks to cache", zap.Error(err))
		} else {
			if err = r.pool.Expire(ctx, key, r.pool.TTL()).Err(); err != nil {
				r.log.Error("set ttl for tasks cache", zap.Error(err))
			}
		}
	}

	return tasks, nil
}
