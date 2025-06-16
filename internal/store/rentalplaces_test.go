package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/testhelpers"
	"github.com/stretchr/testify/suite"
)

type RentalPlacesTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *RentalPlaceStore
	ctx         context.Context
}

func (suite *RentalPlacesTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.repository = NewRentalPlaceStore(suite.pgContainer.DB)
}

func TestRentalPlacesTestSuite(t *testing.T) {
	suite.Run(t, new(RentalPlacesTestSuite))
}

func (suite *RentalPlacesTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *RentalPlacesTestSuite) TestGetRentalPlaceByID() {
	suite.T().Run("it should return an existing rental place", func(t *testing.T) {
		rentalPlace, err := suite.repository.GetRentalPlaceByID(suite.ctx, 1)
		suite.NoError(err)
		suite.NotNil(rentalPlace)
		suite.Equal(rentalPlace.ID, 1)
	})

	suite.T().Run("it should return nil if the rental place does not exist", func(t *testing.T) {
		rentalPlace, err := suite.repository.GetRentalPlaceByID(suite.ctx, 10000)
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(rentalPlace)
	})
}
