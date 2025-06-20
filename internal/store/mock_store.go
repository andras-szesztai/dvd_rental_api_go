package store

import (
	"context"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type MockUserStore struct {
	RegisterUserFunc func(ctx context.Context, user *User) error
	GetUserByIDFunc  func(ctx context.Context, id int64) (*User, error)
}

func (m *MockUserStore) RegisterUser(ctx context.Context, user *User) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

type MockStaffStore struct {
	GetStaffByEmailFunc func(ctx context.Context, email string) (*Staff, error)
}

func (m *MockStaffStore) GetStaffByEmail(ctx context.Context, email string) (*Staff, error) {
	if m.GetStaffByEmailFunc != nil {
		return m.GetStaffByEmailFunc(ctx, email)
	}
	return nil, nil
}

type MockCustomerStore struct {
	CreateCustomerFunc     func(ctx context.Context, customer *Customer) error
	GetCustomerByEmailFunc func(ctx context.Context, email string) (*Customer, error)
}

func (m *MockCustomerStore) CreateCustomer(ctx context.Context, customer *Customer) error {
	if m.CreateCustomerFunc != nil {
		return m.CreateCustomerFunc(ctx, customer)
	}
	return nil
}

func (m *MockCustomerStore) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	if m.GetCustomerByEmailFunc != nil {
		return m.GetCustomerByEmailFunc(ctx, email)
	}
	return nil, nil
}

type MockRoleStore struct {
	GetRoleByNameFunc func(ctx context.Context, name string) (*Role, error)
	GetRoleByIDFunc   func(ctx context.Context, id int64) (*Role, error)
}

func (m *MockRoleStore) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	if m.GetRoleByNameFunc != nil {
		return m.GetRoleByNameFunc(ctx, name)
	}
	if name == "admin" {
		return &Role{
			ID:   1,
			Name: "admin",
		}, nil
	} else if name == "customer" {
		return &Role{
			ID:   2,
			Name: "customer",
		}, nil
	}
	return nil, nil
}

func (m *MockRoleStore) GetRoleByID(ctx context.Context, id int64) (*Role, error) {
	if m.GetRoleByIDFunc != nil {
		return m.GetRoleByIDFunc(ctx, id)
	}
	if id == 1 {
		return &Role{
			ID:   1,
			Name: "admin",
		}, nil
	}
	if id == 2 {
		return &Role{
			ID:   2,
			Name: "customer",
		}, nil
	}
	return nil, nil
}

type MockRentalStore struct {
	GetRentalFunc       func(ctx context.Context, id int64) (*Rental, error)
	GetMovieRentalsFunc func(ctx context.Context, inventoryID int64) ([]*Rental, error)
}

func (m *MockRentalStore) GetRental(ctx context.Context, id int64) (*Rental, error) {
	if m.GetRentalFunc != nil {
		return m.GetRentalFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRentalStore) GetMovieRentals(ctx context.Context, inventoryID int64) ([]*Rental, error) {
	if m.GetMovieRentalsFunc != nil {
		return m.GetMovieRentalsFunc(ctx, inventoryID)
	}
	return nil, nil
}

type MockMovieStore struct {
	GetMoviesFunc func(ctx context.Context, movieQuery *utils.MovieQuery) ([]*Movie, error)
}

func (m *MockMovieStore) GetMovies(ctx context.Context, movieQuery *utils.MovieQuery) ([]*Movie, error) {
	if m.GetMoviesFunc != nil {
		return m.GetMoviesFunc(ctx, movieQuery)
	}
	return nil, nil
}

type MockInventoryStore struct {
	GetMovieInventoryFunc func(ctx context.Context, filmID int) ([]*Inventory, error)
}

func (m *MockInventoryStore) GetMovieInventory(ctx context.Context, filmID int) ([]*Inventory, error) {
	if m.GetMovieInventoryFunc != nil {
		return m.GetMovieInventoryFunc(ctx, filmID)
	}
	return nil, nil
}

func NewMockStore() *Store {
	return &Store{
		Users:     &MockUserStore{},
		Staff:     &MockStaffStore{},
		Customers: &MockCustomerStore{},
		Roles:     &MockRoleStore{},
		Rentals:   &MockRentalStore{},
		Movies:    &MockMovieStore{},
		Inventory: &MockInventoryStore{},
	}
}
