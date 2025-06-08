package store

import (
	"context"
	"database/sql"
)

type RentalPlaceStore struct {
	db *sql.DB
}

func NewRentalPlaceStore(db *sql.DB) *RentalPlaceStore {
	return &RentalPlaceStore{db: db}
}

type RentalPlace struct {
	ID int `json:"id"`
}

func (s *RentalPlaceStore) GetRentalPlaceByID(ctx context.Context, id int) (*RentalPlace, error) {
	query := `
		SELECT store_id
		FROM store
		WHERE store_id = $1
	`

	row := s.db.QueryRowContext(ctx, query, id)

	var rentalPlace RentalPlace
	err := row.Scan(&rentalPlace.ID)
	if err != nil {
		return nil, err
	}

	return &rentalPlace, nil
}
