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
	Rentals interface {
		GetRental(ctx context.Context, id int64) (*Rental, error)
	}
	Users interface {
		CreateAdminUser(ctx context.Context, user *User) error
		GetAdminUserByEmail(ctx context.Context, email string) (*User, error)
		GetAdminUserByID(ctx context.Context, id int64) (*User, error)
	}
	RentalPlaces interface {
		GetRentalPlaceByID(ctx context.Context, id int64) (*RentalPlace, error)
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Rentals:      NewRentalStore(db),
		Users:        NewUserStore(db),
		RentalPlaces: NewRentalPlaceStore(db),
	}
}
