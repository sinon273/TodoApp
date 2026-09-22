package main

import (
	core_logger "TodoApp/internal/core/logger"
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	core_postgres_pool_pgx "TodoApp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "TodoApp/internal/core/transport/http/middleware"
	core_http_server "TodoApp/internal/core/transport/http/server"
	tasks_postgres_repository "TodoApp/internal/features/tasks/repository/postgres"
	task_service "TodoApp/internal/features/tasks/service"
	tasks_transport "TodoApp/internal/features/tasks/transport/http"
	users_postgres_repository "TodoApp/internal/features/users/repository/postgres"
	users_service "TodoApp/internal/features/users/service"
	users_transport_http "TodoApp/internal/features/users/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
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

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTaskRepository(pool)
	tasksService := task_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
