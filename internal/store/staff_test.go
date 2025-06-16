package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/testhelpers"
	"github.com/stretchr/testify/suite"
)

type StaffTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *StaffStore
	ctx         context.Context
}

func (suite *StaffTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.repository = NewStaffStore(suite.pgContainer.DB)
}

func TestStaffTestSuite(t *testing.T) {
	suite.Run(t, new(StaffTestSuite))
}

func (suite *StaffTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *StaffTestSuite) TestGetStaffByEmail() {
	suite.T().Run("it should return an existing staff member", func(t *testing.T) {
		staff, err := suite.repository.GetStaffByEmail(suite.ctx, "Mike.Hillyer@sakilastaff.com")
		suite.NoError(err)
		suite.NotNil(staff)
		suite.Equal(staff.ID, 1)
		suite.Nil(staff.UserID)
	})

	suite.T().Run("it should return nil if the staff member does not exist", func(t *testing.T) {
		staff, err := suite.repository.GetStaffByEmail(suite.ctx, "nonexistent@example.com")
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(staff)
	})
}
