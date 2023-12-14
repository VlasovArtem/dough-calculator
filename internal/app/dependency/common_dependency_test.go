package dependency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

func TestCommonDependencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CommonDependencyServiceTestSuite))
}

type CommonDependencyServiceTestSuite struct {
	test.GoMockTestSuite

	actuatorHandler *mocks.MockActuatorHandler
	configManager   *mocks.MockConfigManager
	mongoDBService  *mocks.MockMongoDBService

	target domain.CommonDependencyService
}

func (suite *CommonDependencyServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.actuatorHandler = mocks.NewMockActuatorHandler(suite.MockCtrl)
	suite.configManager = mocks.NewMockConfigManager(suite.MockCtrl)
	suite.mongoDBService = mocks.NewMockMongoDBService(suite.MockCtrl)

	suite.target = newCommonDependencyService(func() domain.ConfigManager {
		return suite.configManager
	}, func() domain.ActuatorHandler {
		return suite.actuatorHandler
	}, func(config config.Database) (domain.MongoDBService, error) {
		return suite.mongoDBService, nil
	})
}

func (suite *CommonDependencyServiceTestSuite) TestInitialize() {
	suite.configManager.EXPECT().
		ParseConfig().
		Return(nil)
	suite.configManager.EXPECT().
		GetConfig().
		Return(config.Config{})

	err := suite.target.Initialize(context.Background())

	suite.NoError(err)

	suite.Equal(suite.configManager, suite.target.ConfigManager())
	suite.Equal(suite.actuatorHandler, suite.target.Actuator())
	suite.Equal(suite.mongoDBService, suite.target.MongoDBService())
}

func (suite *CommonDependencyServiceTestSuite) TestInitialize_WithErrorOnParseConfig() {
	suite.configManager.EXPECT().
		ParseConfig().
		Return(assert.AnError)

	err := suite.target.Initialize(context.Background())

	suite.ErrorContains(err, "failed to parse config")

	suite.Nil(suite.target.ConfigManager())
}

func (suite *CommonDependencyServiceTestSuite) TestInitialize_WithErrorOnMongoDBServiceCreator() {
	suite.configManager.EXPECT().
		ParseConfig().
		Return(nil)
	suite.configManager.EXPECT().
		GetConfig().
		Return(config.Config{})

	suite.target = newCommonDependencyService(func() domain.ConfigManager {
		return suite.configManager
	}, func() domain.ActuatorHandler {
		return suite.actuatorHandler
	}, func(config config.Database) (domain.MongoDBService, error) {
		return nil, assert.AnError
	})

	err := suite.target.Initialize(context.Background())

	suite.ErrorContains(err, "failed to create mongodb service")

	suite.Nil(suite.target.MongoDBService())
}

func TestNewCommonDependencyService(t *testing.T) {
	service := NewCommonDependencyService().(*commonDependencyService)
	assert.NotNil(t, service)
	assert.NotNil(t, service.configManagerCreator)
	assert.Nil(t, service.configManager)
	assert.NotNil(t, service.actuatorHandlerCreator)
	assert.Nil(t, service.actuatorHandler)
	assert.NotNil(t, service.mongoDBServiceCreator)
	assert.Nil(t, service.mongoDBService)
}
