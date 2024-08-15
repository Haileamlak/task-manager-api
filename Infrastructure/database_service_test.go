package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseServiceTestSuite struct {
	suite.Suite
	dbService DatabaseService
}

func (suite *DatabaseServiceTestSuite) SetupTest() {
	suite.dbService = NewDatabase()
}

func TestDatabaseServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseServiceTestSuite))
}

func (suite *DatabaseServiceTestSuite) TestConnect() {
	db := suite.dbService.Connect("mongodb://localhost:27017")

	assert.NotNil(suite.T(), db)
	assert.Equal(suite.T(), "task_manager", db.Name())
}