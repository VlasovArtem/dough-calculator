package dependency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

type FlourDependencyServiceTestSuite struct {
	test.GoMockTestSuite

	configManager  *mocks.MockConfigManager
	mongoDBService *mocks.MockMongoDBService
	repository     *mocks.MockFlourRepository
	service        *mocks.MockFlourService
	handler        *mocks.MockFlourHandler

	target domain.FlourDependencyService
}

func (suite *FlourDependencyServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.mongoDBService = mocks.NewMockMongoDBService(suite.MockCtrl)
	suite.configManager = mocks.NewMockConfigManager(suite.MockCtrl)
	suite.repository = mocks.NewMockFlourRepository(suite.MockCtrl)
	suite.service = mocks.NewMockFlourService(suite.MockCtrl)
	suite.handler = mocks.NewMockFlourHandler(suite.MockCtrl)

	suite.target = newFlourDependencyService(
		func(_ domain.MongoDBService) (domain.FlourRepository, error) {
			return suite.repository, nil
		},
		func(_ domain.FlourRepository) (domain.FlourService, error) {
			return suite.service, nil
		},
		func(_ domain.FlourService) (domain.FlourHandler, error) {
			return suite.handler, nil
		},
	)
}

func (suite *FlourDependencyServiceTestSuite) TestInitialize() {
	ctx := context.WithValue(context.Background(), "mongoDBService", suite.mongoDBService)

	err := suite.target.Initialize(ctx)

	suite.NoError(err)
	suite.Equal(suite.repository, suite.target.Repository())
	suite.Equal(suite.service, suite.target.Service())
	suite.Equal(suite.handler, suite.target.Router())
}

func (suite *FlourDependencyServiceTestSuite) TestInitialize_MongoDBServiceNil() {
	ctx := context.Background()

	err := suite.target.Initialize(ctx)

	suite.ErrorContains(err, "failed to get mongoDBService from context")
	suite.Nil(suite.target.Repository())
	suite.Nil(suite.target.Service())
	suite.Nil(suite.target.Router())
}

func (suite *FlourDependencyServiceTestSuite) TestInitialize_WithError() {
	baseService := flourDependencyService{
		repositoryCreator: func(_ domain.MongoDBService) (domain.FlourRepository, error) {
			return suite.repository, nil
		},
		serviceCreator: func(_ domain.FlourRepository) (domain.FlourService, error) {
			return suite.service, nil
		},
		handlerCreator: func(_ domain.FlourService) (domain.FlourHandler, error) {
			return suite.handler, nil
		},
	}

	tests := []struct {
		name             string
		serviceCreator   func(service flourDependencyService) domain.FlourDependencyService
		expectedErrorMsg string
	}{
		{
			name: "repositoryCreator",
			serviceCreator: func(service flourDependencyService) domain.FlourDependencyService {
				service.repositoryCreator = func(_ domain.MongoDBService) (domain.FlourRepository, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create repository",
		},
		{
			name: "serviceCreator",
			serviceCreator: func(service flourDependencyService) domain.FlourDependencyService {
				service.serviceCreator = func(_ domain.FlourRepository) (domain.FlourService, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create service",
		},
		{
			name: "handlerCreator",
			serviceCreator: func(service flourDependencyService) domain.FlourDependencyService {
				service.handlerCreator = func(_ domain.FlourService) (domain.FlourHandler, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create handler",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ctx := context.WithValue(context.Background(), "mongoDBService", suite.mongoDBService)

			service := tt.serviceCreator(baseService)

			err := service.Initialize(ctx)

			suite.ErrorContains(err, tt.expectedErrorMsg)
			suite.Nil(service.Repository())
			suite.Nil(service.Service())
			suite.Nil(service.Router())
		})
	}
}

func (suite *FlourDependencyServiceTestSuite) TestRepository() {
	target := &flourDependencyService{
		repository: suite.repository,
	}

	suite.Equal(suite.repository, target.Repository())
}

func (suite *FlourDependencyServiceTestSuite) TestService() {
	target := &flourDependencyService{
		service: suite.service,
	}

	suite.Equal(suite.service, target.Service())
}

func (suite *FlourDependencyServiceTestSuite) TestRouter() {
	target := &flourDependencyService{
		handler: suite.handler,
	}

	suite.Equal(suite.handler, target.Router())
}

func (suite *FlourDependencyServiceTestSuite) TestNewFlourDependencyService() {
	target := NewFlourDependencyService().(*flourDependencyService)

	suite.NotNil(target)
	suite.NotNil(target.repositoryCreator)
	suite.NotNil(target.serviceCreator)
	suite.NotNil(target.handlerCreator)
	suite.Nil(target.repository)
	suite.Nil(target.service)
	suite.Nil(target.handler)
}

func TestFlourDependencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(FlourDependencyServiceTestSuite))
}
