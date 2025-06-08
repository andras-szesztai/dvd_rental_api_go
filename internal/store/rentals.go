package store

import (
	"context"
	"database/sql"
)

type RentalStore struct {
	db *sql.DB
}

func NewRentalStore(db *sql.DB) *RentalStore {
	return &RentalStore{db: db}
}

type Rental struct {
	ID int `json:"id"`
}

func (s *RentalStore) GetRental(ctx context.Context, id int64) (*Rental, error) {
	query := `
		SELECT rental_id
		FROM rental
		WHERE rental_id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)

	var rental Rental
	err := row.Scan(&rental.ID)
	if err != nil {
		return nil, err
	}

	return &rental, nil
}
