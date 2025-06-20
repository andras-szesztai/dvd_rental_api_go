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
	ID          int          `json:"id"`
	RentalDate  time.Time    `json:"rental_date"`
	ReturnDate  sql.NullTime `json:"return_date"`
	InventoryID int          `json:"inventory_id"`
}

func (s *RentalStore) GetRental(ctx context.Context, id int64) (*Rental, error) {
	query := `
		SELECT rental_id,rental_date
		FROM rental
		WHERE rental_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	var rental Rental
	err := row.Scan(&rental.ID, &rental.RentalDate)
	if err != nil {
		return nil, err
	}

	return &rental, nil
}

func (s *RentalStore) GetMovieRentals(ctx context.Context, inventoryID int64) ([]*Rental, error) {
	query := `
		SELECT rental_id, rental_date, return_date, inventory_id
		FROM rental
		WHERE inventory_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, inventoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rentals := []*Rental{}
	for rows.Next() {
		var rental Rental
		if err := rows.Scan(&rental.ID, &rental.RentalDate, &rental.ReturnDate, &rental.InventoryID); err != nil {
			return nil, err
		}
		rentals = append(rentals, &rental)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rentals, nil
}
