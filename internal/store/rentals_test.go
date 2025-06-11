package store

import (
	"context"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/testhelpers"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type RentalsTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *RentalStore
	ctx         context.Context
}

func (suite *RentalsTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.repository = NewRentalStore(suite.pgContainer.DB)
}

func TestRentalsTestSuite(t *testing.T) {
	suite.Run(t, new(RentalsTestSuite))
}

func (suite *RentalsTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *RentalsTestSuite) TestGetRental() {
	rental, err := suite.repository.GetRental(suite.ctx, 1)

	suite.NoError(err)
	suite.NotNil(rental)
	suite.Equal(1, rental.ID)
	suite.Equal("2005-05-24 22:53:30 +0000 +0000", rental.RentalDate.String())

	rental, err = suite.repository.GetRental(suite.ctx, 2)
	suite.NoError(err)
	suite.NotNil(rental)
	suite.Equal(2, rental.ID)
	suite.Equal("2005-05-24 22:54:33 +0000 +0000", rental.RentalDate.String())
}
