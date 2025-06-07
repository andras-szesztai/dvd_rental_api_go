package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Rental interface {
		GetRental(ctx context.Context, id int) (*Rental, error)
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Rental: NewRentalStore(db),
	}
}
