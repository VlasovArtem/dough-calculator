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

type SourdoughRecipeDependencyServiceTestSuite struct {
	test.GoMockTestSuite

	configManager  *mocks.MockConfigManager
	mongoDBService *mocks.MockMongoDBService
	repository     *mocks.MockSourdoughRecipeRepository
	service        *mocks.MockSourdoughRecipeService
	handler        *mocks.MockSourdoughRecipeHandler

	target domain.SourdoughRecipeDependencyService
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.mongoDBService = mocks.NewMockMongoDBService(suite.MockCtrl)
	suite.configManager = mocks.NewMockConfigManager(suite.MockCtrl)
	suite.repository = mocks.NewMockSourdoughRecipeRepository(suite.MockCtrl)
	suite.service = mocks.NewMockSourdoughRecipeService(suite.MockCtrl)
	suite.handler = mocks.NewMockSourdoughRecipeHandler(suite.MockCtrl)

	suite.target = newSourdoughRecipeDependencyService(
		func(_ domain.MongoDBService) (domain.SourdoughRecipeRepository, error) {
			return suite.repository, nil
		},
		func(_ domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error) {
			return suite.service, nil
		},
		func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error) {
			return suite.handler, nil
		},
	)
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestInitialize() {
	ctx := context.WithValue(context.Background(), "mongoDBService", suite.mongoDBService)

	err := suite.target.Initialize(ctx)

	suite.NoError(err)
	suite.Equal(suite.repository, suite.target.Repository())
	suite.Equal(suite.service, suite.target.Service())
	suite.Equal(suite.handler, suite.target.Router())
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestInitialize_MongoDBServiceNil() {
	ctx := context.Background()

	err := suite.target.Initialize(ctx)

	suite.ErrorContains(err, "failed to get mongoDBService from context")
	suite.Nil(suite.target.Repository())
	suite.Nil(suite.target.Service())
	suite.Nil(suite.target.Router())
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestInitialize_WithError() {
	baseService := sourdoughRecipeDependencyService{
		repositoryCreator: func(_ domain.MongoDBService) (domain.SourdoughRecipeRepository, error) {
			return suite.repository, nil
		},
		serviceCreator: func(_ domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error) {
			return suite.service, nil
		},
		handlerCreator: func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error) {
			return suite.handler, nil
		},
	}

	tests := []struct {
		name             string
		serviceCreator   func(service sourdoughRecipeDependencyService) domain.SourdoughRecipeDependencyService
		expectedErrorMsg string
	}{
		{
			name: "repositoryCreator",
			serviceCreator: func(service sourdoughRecipeDependencyService) domain.SourdoughRecipeDependencyService {
				service.repositoryCreator = func(_ domain.MongoDBService) (domain.SourdoughRecipeRepository, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create repository",
		},
		{
			name: "serviceCreator",
			serviceCreator: func(service sourdoughRecipeDependencyService) domain.SourdoughRecipeDependencyService {
				service.serviceCreator = func(_ domain.SourdoughRecipeRepository) (domain.SourdoughRecipeService, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create service",
		},
		{
			name: "handlerCreator",
			serviceCreator: func(service sourdoughRecipeDependencyService) domain.SourdoughRecipeDependencyService {
				service.handlerCreator = func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeHandler, error) {
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

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestRepository() {
	target := &sourdoughRecipeDependencyService{
		repository: suite.repository,
	}

	suite.Equal(suite.repository, target.Repository())
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestService() {
	target := &sourdoughRecipeDependencyService{
		service: suite.service,
	}

	suite.Equal(suite.service, target.Service())
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestRouter() {
	target := &sourdoughRecipeDependencyService{
		handler: suite.handler,
	}

	suite.Equal(suite.handler, target.Router())
}

func (suite *SourdoughRecipeDependencyServiceTestSuite) TestNewSourdoughRecipeDependencyService() {
	target := NewSourdoughRecipeDependencyService().(*sourdoughRecipeDependencyService)

	suite.NotNil(target)
	suite.NotNil(target.repositoryCreator)
	suite.NotNil(target.serviceCreator)
	suite.NotNil(target.handlerCreator)
	suite.Nil(target.repository)
	suite.Nil(target.service)
	suite.Nil(target.handler)
}

func TestSourdoughRecipeDependencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeDependencyServiceTestSuite))
}
