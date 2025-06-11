package store

import (
	"context"
	"database/sql"
	"time"
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

func (s *RentalPlaceStore) GetRentalPlaceByID(ctx context.Context, id int64) (*RentalPlace, error) {
	query := `
		SELECT store_id
		FROM store
		WHERE store_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	var rentalPlace RentalPlace
	err := row.Scan(&rentalPlace.ID)
	if err != nil {
		return nil, err
	}

	return &rentalPlace, nil
}
