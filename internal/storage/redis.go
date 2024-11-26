package storage

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) Storage {
	return &redisStorage{client: client}
}

var ErrNotFound = errors.New("not found")

func (s *redisStorage) Get(ctx context.Context, key string) (int, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	// Converter o valor para inteiro
	tokens, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}

func (s *redisStorage) Increment(ctx context.Context, key string, value int) error {
	return s.client.IncrBy(ctx, key, int64(value)).Err()
}

func (s *redisStorage) Decrement(ctx context.Context, key string) error {
	return s.client.Decr(ctx, key).Err()
}

func (s *redisStorage) HSet(ctx context.Context, key string, value interface{}) error {
	return s.client.HSet(ctx, key, value).Err()
}

func (s *redisStorage) HGetAll(ctx context.Context, key string) (interface{}, error) {
	values, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return values, nil
}
