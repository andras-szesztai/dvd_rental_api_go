package store

import (
	"context"
	"database/sql"
)

type CustomerStore struct {
	db *sql.DB
}

func NewCustomerStore(db *sql.DB) *CustomerStore {
	return &CustomerStore{db: db}
}

type Customer struct {
	ID     int  `json:"id"`
	UserID *int `json:"user_id"`
}

func (s *CustomerStore) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	query := `
		SELECT customer_id, user_id
		FROM customer
		WHERE email = $1
	`

	row := s.db.QueryRowContext(ctx, query, email)

	var customer Customer
	var userID sql.NullInt64
	err := row.Scan(&customer.ID, &userID)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		userIDInt := int(userID.Int64)
		customer.UserID = &userIDInt
	}

	return &customer, nil
}
