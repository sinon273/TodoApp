package task_cached_repository

import (
	core_logger "TodoApp/internal/core/logger"
	core_redis_pool "TodoApp/internal/core/repository/redis/pool"
	task_ports_repository "TodoApp/internal/features/tasks/ports"
)

type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository task_ports_repository.TaskRepository
	log            *core_logger.Logger
}

func NewCachedRepository(
	pool core_redis_pool.Pool,
	mainRepository task_ports_repository.TaskRepository,
	log *core_logger.Logger,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
		log:            log,
	}
}
