package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type Store struct {
	Rental interface {
		GetRental(ctx context.Context, id int) (*Rental, error)
	}
	User interface {
		CreateAdminUser(ctx context.Context, user *AdminUser) error
	}
	RentalPlace interface {
		GetRentalPlaceByID(ctx context.Context, id int) (*RentalPlace, error)
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Rental:      NewRentalStore(db),
		User:        NewUserStore(db),
		RentalPlace: NewRentalPlaceStore(db),
	}
}
