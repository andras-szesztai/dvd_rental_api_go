package store

import (
	"context"
	"database/sql"
	"time"
)

type RoleStore struct {
	db *sql.DB
}

func NewRoleStore(db *sql.DB) *RoleStore {
	return &RoleStore{db: db}
}

type Role struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (s *RoleStore) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, level
		FROM roles
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, name)

	var role Role
	err := row.Scan(&role.ID, &role.Name, &role.Level)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleStore) GetRoleByID(ctx context.Context, id int64) (*Role, error) {
	query := `
		SELECT id, name, level
		FROM roles
		WHERE id = $1
	`

	row := s.db.QueryRowContext(ctx, query, id)

	var role Role
	err := row.Scan(&role.ID, &role.Name, &role.Level)
	if err != nil {
		return nil, err
	}

	return &role, nil
}
