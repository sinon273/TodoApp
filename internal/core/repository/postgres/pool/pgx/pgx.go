package core_postgres_pool_pgx

import (
	core_postgres_pool "TodoApp/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config core_postgres_pool.Config,
) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgxpool: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	return &row{Row: p.Pool.QueryRow(ctx, sql, args...)}
}

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}
	return &rowsWrapper{Rows: rows}, nil
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return &commandTag{CommandTag: tag}, nil
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return core_postgres_pool.ErrViolatesForeignKey
	}
	return err
}

type row struct {
	pgx.Row
}

type rowsWrapper struct {
	pgx.Rows
}

func (r *rowsWrapper) Err() error {
	if err := r.Rows.Err(); err != nil {
		return mapErrors(err)
	}
	return nil
}

func (r *row) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

type commandTag struct {
	pgconn.CommandTag
}
