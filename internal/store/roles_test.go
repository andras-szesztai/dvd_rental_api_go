package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"github.com/stretchr/testify/suite"
)

type RolesTestSuite struct {
	suite.Suite
	pgContainer *utils.PostgresContainer
	repository  *RoleStore
	ctx         context.Context
}

func (suite *RolesTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := utils.CreatePostgresContainer()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.repository = NewRoleStore(suite.pgContainer.DB)
}

func TestRolesTestSuite(t *testing.T) {
	suite.Run(t, new(RolesTestSuite))
}

func (suite *RolesTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *RolesTestSuite) TestGetRoleByName() {
	suite.T().Run("it should return an existing role", func(t *testing.T) {
		role, err := suite.repository.GetRoleByName(suite.ctx, "admin")
		suite.NoError(err)
		suite.NotNil(role)
		suite.Equal(role.ID, 1)
		suite.Equal(role.Name, "admin")
		suite.Equal(role.Level, 10)
	})

	suite.T().Run("it should return nil if the role does not exist", func(t *testing.T) {
		role, err := suite.repository.GetRoleByName(suite.ctx, "nonexistent")
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(role)
	})
}

func (suite *RolesTestSuite) TestGetRoleByID() {
	suite.T().Run("it should return an existing role", func(t *testing.T) {
		role, err := suite.repository.GetRoleByID(suite.ctx, 1)
		suite.NoError(err)
		suite.NotNil(role)
		suite.Equal(role.ID, 1)
		suite.Equal(role.Name, "admin")
		suite.Equal(role.Level, 10)
	})

	suite.T().Run("it should return nil if the role does not exist", func(t *testing.T) {
		role, err := suite.repository.GetRoleByID(suite.ctx, 10000)
		suite.True(errors.Is(err, sql.ErrNoRows))
		suite.Nil(role)
	})
}
