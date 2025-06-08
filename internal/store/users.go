package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

type AdminUser struct {
	ID         int            `json:"id"`
	AddressID  int            `json:"address_id"`
	StoreID    int            `json:"store_id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	Password   utils.Password `json:"-"`
	LastUpdate time.Time      `json:"last_update"`
}

func (s *UserStore) CreateAdminUser(ctx context.Context, user *AdminUser) error {
	query := `
		INSERT INTO staff (store_id, first_name, last_name, email, username, password)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query, user.StoreID, user.FirstName, user.LastName, user.Email, user.Username, user.Password.Hash)
	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"staff_email_key\"":
			return ErrEmailAlreadyExists
		case "pq: duplicate key value violates unique constraint \"staff_username_key\"":
			return ErrUsernameAlreadyExists
		default:
			return err
		}
	}

	return nil
}
