package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type Store struct {
	Users interface {
		RegisterUser(ctx context.Context, user *User) error
		GetUserByID(ctx context.Context, id int64) (*User, error)
	}
	Staff interface {
		GetStaffByEmail(ctx context.Context, email string) (*Staff, error)
	}
	Customers interface {
		CreateCustomer(ctx context.Context, customer *Customer) error
		GetCustomerByEmail(ctx context.Context, email string) (*Customer, error)
	}
	Roles interface {
		GetRoleByName(ctx context.Context, name string) (*Role, error)
		GetRoleByID(ctx context.Context, id int64) (*Role, error)
	}
	Movies interface {
		GetMovies(ctx context.Context, movieQuery *utils.MovieQuery) ([]*Movie, error)
	}
	Rentals interface {
		GetRental(ctx context.Context, id int64) (*Rental, error)
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
		Staff:        NewStaffStore(db),
		Customers:    NewCustomerStore(db),
		Roles:        NewRoleStore(db),
		Movies:       NewMovieStore(db),
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
