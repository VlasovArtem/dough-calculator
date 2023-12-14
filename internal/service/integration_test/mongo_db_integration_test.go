//go:build integration && docker

package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/service"
	"dough-calculator/internal/test"
)

const (
	testDb         = "test_db"
	testCollection = "test_collection"
)

func TestMongoDBTestSuite(t *testing.T) {
	suite.Run(t, &MongoDBTestSuite{
		MongoDBDockerIntegrationTestSuite: test.NewMongoDBDockerIntegrationTestSuite(dockerStarter),
	})
}

type MongoDBTestSuite struct {
	test.MongoDBDockerIntegrationTestSuite

	service domain.MongoDBService
}

func (suite *MongoDBTestSuite) SetupSuite() {
	suite.MongoDBDockerIntegrationTestSuite.SetupSuite()

	suite.service = test.Must(func() (domain.MongoDBService, error) {
		return service.NewMongoDBService(suite.GetConfig())
	})
}

func (suite *MongoDBTestSuite) AfterTest(suiteName, testName string) {
	collection, err := suite.service.GetCollection(testDb, testCollection)
	suite.Require().NoError(err)

	err = collection.Drop(context.Background())
}

func (suite *MongoDBTestSuite) TestGetClient() {
	newClient, err := suite.service.GetClient()
	suite.Require().NoError(err)
	suite.Require().NotNil(newClient)

	err = newClient.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	existingClient, err := suite.service.GetClient()

	suite.Require().NoError(err)
	suite.Require().NotNil(existingClient)

	suite.Equal(*newClient, *existingClient)
}

func (suite *MongoDBTestSuite) TestGetDatabase() {
	newDatabase, err := suite.service.GetDatabase(testDb)
	suite.Require().NoError(err)
	suite.Require().NotNil(newDatabase)

	existingDatabase, err := suite.service.GetDatabase(testDb)

	suite.Require().NoError(err)
	suite.Require().NotNil(existingDatabase)

	suite.Equal(*newDatabase, *existingDatabase)
}

func (suite *MongoDBTestSuite) TestGetCollection() {
	newCollection, err := suite.service.GetCollection(testDb, testCollection)
	suite.Require().NoError(err)
	suite.Require().NotNil(newCollection)

	existingCollection, err := suite.service.GetCollection(testDb, testCollection)

	suite.Require().NoError(err)
	suite.Require().NotNil(existingCollection)

	suite.Equal(*newCollection, *existingCollection)
}
