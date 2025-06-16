package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/testhelpers"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type CustomersTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *CustomerStore
	ctx         context.Context
}

func (suite *CustomersTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.repository = NewCustomerStore(suite.pgContainer.DB)

}

func TestCustomersTestSuite(t *testing.T) {
	suite.Run(t, new(CustomersTestSuite))
}

func (suite *CustomersTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *CustomersTestSuite) TestGetCustomerByEmail() {
	suite.T().Run("it should return an existing customer", func(t *testing.T) {
		customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "lisa.anderson@sakilacustomer.org")

		suite.NoError(err)
		suite.NotNil(customer)
		suite.Equal(11, customer.ID)
		suite.Nil(customer.UserID)
		suite.Empty(customer.FirstName)
		suite.Empty(customer.LastName)
		suite.Empty(customer.Email)
		suite.Empty(customer.Phone)
		suite.Empty(customer.StoreID)
		suite.Empty(customer.Address)
	})

	suite.T().Run("it should return nil if the customer does not exist", func(t *testing.T) {
		customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "non-existing@example.com")
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(customer)
	})
}

func (suite *CustomersTestSuite) TestCreatCustomer() {
	suite.T().Run("it should create a new customer", func(t *testing.T) {
		customer := &Customer{
			StoreID:   1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "test@example.com",
		}

		// First check that the customer doesn't exist
		existingCustomer, err := suite.repository.GetCustomerByEmail(suite.ctx, "test@example.com")
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(existingCustomer)

		// Create the new customer
		err = suite.repository.CreateCustomer(suite.ctx, customer)
		suite.NoError(err)

		// Verify the customer was created
		createdCustomer, err := suite.repository.GetCustomerByEmail(suite.ctx, "test@example.com")
		suite.NoError(err)
		suite.NotNil(createdCustomer)
		suite.Equal(createdCustomer.ID, 10000)
		suite.Nil(createdCustomer.UserID)
	})
}
