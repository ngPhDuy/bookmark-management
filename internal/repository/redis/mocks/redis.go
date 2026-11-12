package mock

import (
	"context"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
)

type mockRedisRepository struct {
	Store map[string]string

	GetErr error
	SetErr error

	CountSetCall int
	LastExp      int
}

func NewMockRedisRepository() *mockRedisRepository {
	return &mockRedisRepository{
		Store: map[string]string{},
	}
}

func (m *mockRedisRepository) SetIfNotExist(ctx context.Context, key string, value string, exp int) (string, error) {
	m.CountSetCall++

	if m.SetErr != nil {
		return "", m.SetErr
	}

	if _, exist := m.Store[key]; exist {
		return "", domain.ErrKeyExists
	}

	m.Store[key] = value
	m.LastExp = exp
	return key, nil
}

func (m *mockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	if m.GetErr != nil {
		return "", m.GetErr
	}

	value, exist := m.Store[key]

	if !exist {
		return "", domain.ErrKeyNotFound
	}

	return value, nil
}
