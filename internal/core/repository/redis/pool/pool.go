package core_redis_pool

import (
	"context"
	"errors"
	"time"
)

var NotFound = errors.New("not found")

type StringCmd interface {
	Bytes() ([]byte, error)
	Result() (string, error)
}

type StatusCmd interface {
	Err() error
}

type IntCmd interface {
	Result() (int64, error)
	Err() error
}

type BoolCmd interface {
	Result() (bool, error)
	Err() error
}

type Pool interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, ttl time.Duration) StatusCmd
	Del(ctx context.Context, keys ...string) IntCmd
	HGet(ctx context.Context, key string, field string) StringCmd
	HSet(ctx context.Context, key string, values ...any) IntCmd
	Expire(ctx context.Context, key string, ttl time.Duration) BoolCmd
	Close() error
	TTL() time.Duration
}
