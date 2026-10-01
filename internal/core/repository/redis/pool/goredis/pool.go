package core_goredis_pool

import (
	core_redis_pool "TodoApp/internal/core/repository/redis/pool"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPool(ctx context.Context, cfg Config) (*Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{client: client, ttl: cfg.TTL}, nil
}

type stringCmd struct {
	*redis.StringCmd
}

func (c *stringCmd) Bytes() ([]byte, error) {
	bytes, err := c.StringCmd.Bytes()
	return bytes, mapError(err)
}

func (c *stringCmd) Result() (string, error) {
	result, err := c.StringCmd.Result()
	return result, mapError(err)
}

type statusCmd struct {
	*redis.StatusCmd
}

func (c *statusCmd) Err() error {
	return mapError(c.StatusCmd.Err())
}

type intCmd struct {
	*redis.IntCmd
}

type boolCmd struct {
	*redis.BoolCmd
}

func (c *boolCmd) Result() (bool, error) {
	result, err := c.BoolCmd.Result()
	return result, mapError(err)
}

func (c *boolCmd) Err() error {
	return mapError(c.BoolCmd.Err())
}

func (c *intCmd) Result() (int64, error) {
	result, err := c.IntCmd.Result()
	return result, mapError(err)
}

func (c *intCmd) Err() error {
	return mapError(c.IntCmd.Err())
}

func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.NotFound
	}
	return err
}

func (p *Pool) Get(ctx context.Context, key string) core_redis_pool.StringCmd {
	return &stringCmd{StringCmd: p.client.Get(ctx, key)}
}

func (p *Pool) HGet(ctx context.Context, key string, field string) core_redis_pool.StringCmd {
	return &stringCmd{StringCmd: p.client.HGet(ctx, key, field)}
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) core_redis_pool.StatusCmd {
	return &statusCmd{StatusCmd: p.client.Set(ctx, key, value, ttl)}
}
func (p *Pool) Del(ctx context.Context, keys ...string) core_redis_pool.IntCmd {
	return &intCmd{IntCmd: p.client.Del(ctx, keys...)}
}

func (p *Pool) HSet(ctx context.Context, key string, values ...any) core_redis_pool.IntCmd {
	return &intCmd{IntCmd: p.client.HSet(ctx, key, values)}
}
func (p *Pool) Close() error {
	return p.client.Close()
}
func (p *Pool) TTL() time.Duration {
	return p.ttl
}

func (p *Pool) Expire(ctx context.Context, key string, ttl time.Duration) core_redis_pool.BoolCmd {
	return &boolCmd{BoolCmd: p.client.Expire(ctx, key, ttl)}
}
