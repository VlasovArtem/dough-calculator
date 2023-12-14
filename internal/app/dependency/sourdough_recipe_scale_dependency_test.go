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

type SourdoughRecipeScaleDependencyServiceTestSuite struct {
	test.GoMockTestSuite

	sourdoughRecipeService *mocks.MockSourdoughRecipeService
	service                *mocks.MockSourdoughRecipeScaleService
	handler                *mocks.MockSourdoughRecipeScaleHandler

	target domain.SourdoughRecipeScaleDependencyService
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.sourdoughRecipeService = mocks.NewMockSourdoughRecipeService(suite.MockCtrl)
	suite.service = mocks.NewMockSourdoughRecipeScaleService(suite.MockCtrl)
	suite.handler = mocks.NewMockSourdoughRecipeScaleHandler(suite.MockCtrl)

	suite.target = newSourdoughRecipeScaleDependencyService(
		func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error) {
			return suite.service, nil
		},
		func(_ domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error) {
			return suite.handler, nil
		},
	)
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestInitialize() {
	ctx := context.WithValue(context.Background(), "sourdoughRecipeService", suite.sourdoughRecipeService)

	err := suite.target.Initialize(ctx)

	suite.NoError(err)
	suite.Equal(suite.service, suite.target.Service())
	suite.Equal(suite.handler, suite.target.Router())
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestInitialize_SourdoughRecipeServiceNil() {
	ctx := context.Background()

	err := suite.target.Initialize(ctx)

	suite.ErrorContains(err, "failed to get sourdoughRecipeService from context")
	suite.Nil(suite.target.Service())
	suite.Nil(suite.target.Router())
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestInitialize_WithError() {
	baseService := sourdoughRecipeScaleDependencyService{
		serviceCreator: func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error) {
			return suite.service, nil
		},
		handlerCreator: func(_ domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error) {
			return suite.handler, nil
		},
	}

	tests := []struct {
		name             string
		serviceCreator   func(service sourdoughRecipeScaleDependencyService) domain.SourdoughRecipeScaleDependencyService
		expectedErrorMsg string
	}{
		{
			name: "serviceCreator",
			serviceCreator: func(service sourdoughRecipeScaleDependencyService) domain.SourdoughRecipeScaleDependencyService {
				service.serviceCreator = func(_ domain.SourdoughRecipeService) (domain.SourdoughRecipeScaleService, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create service",
		},
		{
			name: "handlerCreator",
			serviceCreator: func(service sourdoughRecipeScaleDependencyService) domain.SourdoughRecipeScaleDependencyService {
				service.handlerCreator = func(_ domain.SourdoughRecipeScaleService) (domain.SourdoughRecipeScaleHandler, error) {
					return nil, assert.AnError
				}

				return &service
			},
			expectedErrorMsg: "failed to create handler",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ctx := context.WithValue(context.Background(), "sourdoughRecipeService", suite.sourdoughRecipeService)

			service := tt.serviceCreator(baseService)

			err := service.Initialize(ctx)

			suite.ErrorContains(err, tt.expectedErrorMsg)
			suite.Nil(service.Service())
			suite.Nil(service.Router())
		})
	}
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestService() {
	target := &sourdoughRecipeScaleDependencyService{
		service: suite.service,
	}

	suite.Equal(suite.service, target.Service())
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestRouter() {
	target := &sourdoughRecipeScaleDependencyService{
		handler: suite.handler,
	}

	suite.Equal(suite.handler, target.Router())
}

func (suite *SourdoughRecipeScaleDependencyServiceTestSuite) TestNewSourdoughRecipeScaleDependencyService() {
	target := NewSourdoughRecipeScaleDependencyService().(*sourdoughRecipeScaleDependencyService)

	suite.NotNil(target)
	suite.NotNil(target.serviceCreator)
	suite.NotNil(target.handlerCreator)
	suite.Nil(target.service)
	suite.Nil(target.handler)
}

func TestSourdoughRecipeScaleDependencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeScaleDependencyServiceTestSuite))
}
