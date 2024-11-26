package storage

import (
	"context"
	"errors"
)

type MockStorage struct {
	Data map[string]interface{}
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Data: make(map[string]interface{}),
	}
}

func (m *MockStorage) Get(ctx context.Context, key string) (int, error) {
	if value, ok := m.Data[key].(int); ok {
		return value, nil
	}
	return 0, errors.New("key not found")
}

func (m *MockStorage) Increment(ctx context.Context, key string, value int) error {
	if oldValue, ok := m.Data[key].(int); ok {
		m.Data[key] = oldValue + value
	} else {
		m.Data[key] = value
	}
	return nil
}

func (m *MockStorage) Decrement(ctx context.Context, key string) error {
	if oldValue, ok := m.Data[key].(int); ok {
		m.Data[key] = oldValue - 1
	} else {
		m.Data[key] = -1
	}
	return nil
}

func (m *MockStorage) HGetAll(ctx context.Context, key string) (interface{}, error) {
	if value, ok := m.Data[key].(map[string]interface{}); ok {
		return value, nil
	}
	return nil, errors.New("key not found")
}

func (m *MockStorage) HSet(ctx context.Context, key string, value interface{}) error {
	if _, ok := m.Data[key].(map[string]interface{}); !ok {
		m.Data[key] = make(map[string]interface{})
	}
	m.Data[key].(map[string]interface{})[key] = value
	return nil
}
