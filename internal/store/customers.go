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
	ID        int    `json:"id"`
	UserID    *int   `json:"user_id"`
	StoreID   int64  `json:"store_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
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

func (s *CustomerStore) CreateCustomer(ctx context.Context, customer *Customer) error {
	query := `
		INSERT INTO customer (store_id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(ctx, query, customer.StoreID, customer.FirstName, customer.LastName, customer.Email)
	if err != nil {
		return err
	}

	return nil
}
