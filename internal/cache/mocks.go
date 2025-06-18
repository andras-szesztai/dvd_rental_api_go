package cache

import (
	"context"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
)

func NewMockCache() *Storage {
	return &Storage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct {
	GetFunc    func(ctx context.Context, id int64) (*store.User, error)
	SetFunc    func(ctx context.Context, user *store.User) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockUserCache) Get(ctx context.Context, id int64) (*store.User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserCache) Set(ctx context.Context, user *store.User) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, user)
	}
	return nil
}

func (m *MockUserCache) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
