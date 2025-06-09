package store

import (
	"context"
)

func NewMockStore() *Store {
	return &Store{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) RegisterUser(ctx context.Context, user *User) error {
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}
