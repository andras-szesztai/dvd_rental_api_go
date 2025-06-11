package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type CustomersTestSuite struct {
	suite.Suite
	pgContainer *utils.PostgresContainer
	repository  *CustomerStore
	ctx         context.Context
}

func (suite *CustomersTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := utils.CreatePostgresContainer()
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

// func (suite *CustomersTestSuite) TestCreatCustomer() {
// 	customer := &Customer{
// 		StoreID:   1,
// 		FirstName: "Lisa",
// 		LastName:  "Anderson",
// 		Email:     "lisa.anderson@sakilacustomer.org",
// 	}
// }
