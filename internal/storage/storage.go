package storage

import "context"

type Storage interface {
	Get(ctx context.Context, key string) (int, error)
	Increment(ctx context.Context, key string, value int) error
	Decrement(ctx context.Context, key string) error
	HGetAll(ctx context.Context, key string) (interface{}, error)
	HSet(ctx context.Context, key string, value interface{}) error
}
