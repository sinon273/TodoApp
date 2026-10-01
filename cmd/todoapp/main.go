package main

import (
	core_logger "TodoApp/internal/core/logger"
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	core_postgres_pool_pgx "TodoApp/internal/core/repository/postgres/pool/pgx"
	core_redis_pool_goredis "TodoApp/internal/core/repository/redis/pool/goredis"
	core_http_middleware "TodoApp/internal/core/transport/http/middleware"
	core_http_server "TodoApp/internal/core/transport/http/server"
	statistics_postgres_repository "TodoApp/internal/features/statistics/repository/postgres"
	statistics_service "TodoApp/internal/features/statistics/service"
	statistics_transport_http "TodoApp/internal/features/statistics/transport/http"
	tasks_http "TodoApp/internal/features/tasks/adapters/task_in/http"
	task_cached_repository "TodoApp/internal/features/tasks/adapters/task_out/cached"
	tasks_postgres_repository "TodoApp/internal/features/tasks/adapters/task_out/postgres"
	task_service "TodoApp/internal/features/tasks/service"
	users_postgres_repository "TodoApp/internal/features/users/repository/postgres"
	users_service "TodoApp/internal/features/users/service"
	users_transport_http "TodoApp/internal/features/users/transport/http"
	web_fs_repository "TodoApp/internal/features/web/repository/file_system"
	web_service "TodoApp/internal/features/web/service"
	web_transport_http "TodoApp/internal/features/web/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	_ "TodoApp/docs"
)

func loadTimeZone() *time.Location {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Println("failed to load time zone, fallback to UTC:", err)
		return time.UTC
	}

	return zone
}

// @title       Golang Todo API
// @version     1.0
// @description Todo Application REST-API scheme
// @host 		127.0.0.1:5050
// @BasePath    /api/v1
func main() {
	timeZone := loadTimeZone()
	time.Local = timeZone

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool_pgx.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing redis connection pool")
	redisPool, err := core_redis_pool_goredis.NewPool(
		ctx,
		core_redis_pool_goredis.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init redis connection pool", zap.Error(err))
	}
	defer redisPool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTaskRepository(pool)
	cachedTasksRepository := task_cached_repository.NewCachedRepository(redisPool, tasksRepository, logger)
	tasksService := task_service.NewTasksService(cachedTasksRepository)
	tasksTransportHTTP := tasks_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
