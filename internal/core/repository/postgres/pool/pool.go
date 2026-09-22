package core_postgres_pool

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Close()
	OpTimeout() time.Duration
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type CommandTag interface {
	RowsAffected() int64
}
