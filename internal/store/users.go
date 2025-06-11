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

type User struct {
	ID       int            `json:"id"`
	Email    string         `json:"email"`
	Username string         `json:"username"`
	Role     *Role          `json:"role"`
	Password utils.Password `json:"-"`
}

func (s *UserStore) RegisterUser(ctx context.Context, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		query := `
			INSERT INTO users (username, role_id, password)
			VALUES ($1, $2, $3)
			RETURNING id
		`

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var userID int
		err := tx.QueryRowContext(ctx, query, user.Username, user.Role.ID, user.Password.Hash).Scan(&userID)
		if err != nil {
			return err
		}

		user.ID = userID

		var updateQuery string
		if user.Role.Name == "admin" {
			updateQuery = `
				UPDATE staff SET user_id = $1 WHERE email = $2
			`
		} else {
			updateQuery = `
				UPDATE customer SET user_id = $1 WHERE email = $2
			`
		}

		_, err = tx.ExecContext(ctx, updateQuery, userID, user.Email)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, role_id, password
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Role.ID, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, role_id, password
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	var user User
	user.Role = &Role{}
	err := row.Scan(&user.ID, &user.Username, &user.Role.ID, &user.Password.Hash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
