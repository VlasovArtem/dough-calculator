package test

import (
	"dough-calculator/internal/config"
)

type MongoDBDockerIntegrationTestSuite struct {
	ApplicationTestSuite
	starter MongoDBStarter
}

func (suite *MongoDBDockerIntegrationTestSuite) GetConfig() config.Database {
	dbConfig, err := suite.starter.GetConfig()
	suite.Require().NoError(err)

	return dbConfig
}

func NewMongoDBDockerIntegrationTestSuite(starter MongoDBStarter) MongoDBDockerIntegrationTestSuite {
	return MongoDBDockerIntegrationTestSuite{
		starter: starter,
	}
}
