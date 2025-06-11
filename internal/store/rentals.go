package store

import (
	"context"
	"database/sql"
	"time"
)

type RentalStore struct {
	db *sql.DB
}

func NewRentalStore(db *sql.DB) *RentalStore {
	return &RentalStore{db: db}
}

type Rental struct {
	ID         int       `json:"id"`
	RentalDate time.Time `json:"rental_date"`
}

func (s *RentalStore) GetRental(ctx context.Context, id int64) (*Rental, error) {
	query := `
		SELECT rental_id,rental_date
		FROM rental
		WHERE rental_id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)

	var rental Rental
	err := row.Scan(&rental.ID, &rental.RentalDate)
	if err != nil {
		return nil, err
	}

	return &rental, nil
}
