package dependency

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

type DependencyManagerTestSuite struct {
	test.GoMockTestSuite

	configManager           *mocks.MockConfigManager
	mongoDBService          *mocks.MockMongoDBService
	commonDependencyService *mocks.MockCommonDependencyService

	sourdoughRecipeService           *mocks.MockSourdoughRecipeService
	sourdoughRecipeDependencyService *mocks.MockSourdoughRecipeDependencyService

	sourdoughRecipeScaleDependencyService *mocks.MockSourdoughRecipeScaleDependencyService

	flourDependencyService *mocks.MockFlourDependencyService

	target domain.DependencyManager
}

func (suite *DependencyManagerTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.configManager = mocks.NewMockConfigManager(suite.MockCtrl)
	suite.mongoDBService = mocks.NewMockMongoDBService(suite.MockCtrl)
	suite.commonDependencyService = mocks.NewMockCommonDependencyService(suite.MockCtrl)

	suite.sourdoughRecipeService = mocks.NewMockSourdoughRecipeService(suite.MockCtrl)
	suite.sourdoughRecipeDependencyService = mocks.NewMockSourdoughRecipeDependencyService(suite.MockCtrl)

	suite.sourdoughRecipeScaleDependencyService = mocks.NewMockSourdoughRecipeScaleDependencyService(suite.MockCtrl)

	suite.flourDependencyService = mocks.NewMockFlourDependencyService(suite.MockCtrl)

	suite.target = newDependencyManager(
		suite.commonDependencyService,
		suite.sourdoughRecipeDependencyService,
		suite.sourdoughRecipeScaleDependencyService,
		suite.flourDependencyService,
	)
}

func (suite *DependencyManagerTestSuite) TestInitialize() {
	ctx := context.Background()

	suite.commonDependencyService.EXPECT().Initialize(gomock.Any()).
		DoAndReturn(func(actualCtx context.Context) error {
			suite.Equal(ctx, actualCtx)
			return nil
		})
	suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager)
	suite.commonDependencyService.EXPECT().MongoDBService().Return(suite.mongoDBService)

	suite.sourdoughRecipeDependencyService.EXPECT().Initialize(gomock.Any()).
		DoAndReturn(func(ctx context.Context) error {
			suite.Equal(suite.configManager, ctx.Value("configManager"))
			suite.Equal(suite.mongoDBService, ctx.Value("mongoDBService"))
			return nil
		})
	suite.sourdoughRecipeDependencyService.EXPECT().Service().Return(suite.sourdoughRecipeService)

	suite.sourdoughRecipeScaleDependencyService.EXPECT().Initialize(gomock.Any()).
		DoAndReturn(func(ctx context.Context) error {
			suite.Equal(suite.configManager, ctx.Value("configManager"))
			suite.Equal(suite.mongoDBService, ctx.Value("mongoDBService"))
			suite.Equal(suite.sourdoughRecipeService, ctx.Value("sourdoughRecipeService"))
			return nil
		})

	suite.flourDependencyService.EXPECT().Initialize(gomock.Any()).
		DoAndReturn(func(ctx context.Context) error {
			suite.Equal(suite.mongoDBService, ctx.Value("mongoDBService"))
			return nil
		})

	err := suite.target.Initialize(ctx)

	suite.NoError(err)
	suite.Equal(suite.sourdoughRecipeDependencyService, suite.target.SourdoughRecipe())
	suite.Equal(suite.sourdoughRecipeScaleDependencyService, suite.target.SourdoughRecipeScale())
	suite.Equal(suite.commonDependencyService, suite.target.Common())
	suite.Equal(suite.flourDependencyService, suite.target.Flour())
}

func (suite *DependencyManagerTestSuite) TestInitialize_WithError() {
	tests := []struct {
		name           string
		initializer    func()
		expectedErrMsg string
	}{
		{
			name:           "CommonDependencyService.Initialize() returns error",
			initializer:    func() { suite.commonDependencyService.EXPECT().Initialize(gomock.Any()).Return(assert.AnError) },
			expectedErrMsg: "failed to initialize common dependency service",
		},
		{
			name: "SourdoughRecipeDependencyService.Initialize() returns error",
			initializer: func() {
				suite.commonDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)
				suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager)
				suite.commonDependencyService.EXPECT().MongoDBService().Return(suite.mongoDBService)

				suite.sourdoughRecipeDependencyService.EXPECT().Initialize(gomock.Any()).Return(assert.AnError)
			},
			expectedErrMsg: "failed to initialize sourdough recipe dependency service",
		},
		{
			name: "SourdoughRecipeScaleDependencyService.Initialize() returns error",
			initializer: func() {
				suite.commonDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)
				suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager)
				suite.commonDependencyService.EXPECT().MongoDBService().Return(suite.mongoDBService)

				suite.sourdoughRecipeDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)
				suite.sourdoughRecipeDependencyService.EXPECT().Service().Return(suite.sourdoughRecipeService)

				suite.sourdoughRecipeScaleDependencyService.EXPECT().Initialize(gomock.Any()).Return(assert.AnError)
			},
			expectedErrMsg: "failed to initialize sourdough recipe scale dependency service",
		},
		{
			name: "FlourDependencyService.Initialize() returns error",
			initializer: func() {
				suite.commonDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)
				suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager)
				suite.commonDependencyService.EXPECT().MongoDBService().Return(suite.mongoDBService)

				suite.sourdoughRecipeDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)
				suite.sourdoughRecipeDependencyService.EXPECT().Service().Return(suite.sourdoughRecipeService)

				suite.sourdoughRecipeScaleDependencyService.EXPECT().Initialize(gomock.Any()).Return(nil)

				suite.flourDependencyService.EXPECT().Initialize(gomock.Any()).Return(assert.AnError)
			},
			expectedErrMsg: "failed to initialize flour dependency service",
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			tt.initializer()

			err := suite.target.Initialize(context.Background())

			suite.ErrorContains(err, tt.expectedErrMsg)
		})
	}
}

func (suite *DependencyManagerTestSuite) TestSourdoughRecipe() {
	target := &dependencyManager{
		sourdoughRecipeDependencyService: suite.sourdoughRecipeDependencyService,
	}

	suite.Equal(suite.sourdoughRecipeDependencyService, target.SourdoughRecipe())
}

func (suite *DependencyManagerTestSuite) TestSourdoughRecipeScale() {
	target := &dependencyManager{
		sourdoughRecipeScaleDependencyService: suite.sourdoughRecipeScaleDependencyService,
	}

	suite.Equal(suite.sourdoughRecipeScaleDependencyService, target.SourdoughRecipeScale())
}

func (suite *DependencyManagerTestSuite) TestCommon() {
	target := &dependencyManager{
		commonDependencyService: suite.commonDependencyService,
	}

	suite.Equal(suite.commonDependencyService, target.Common())
}

func (suite *DependencyManagerTestSuite) TestFlour() {
	target := &dependencyManager{
		flourDependencyService: suite.flourDependencyService,
	}

	suite.Equal(suite.flourDependencyService, target.Flour())
}

func (suite *DependencyManagerTestSuite) TestNewDependencyManager() {
	target := NewDependencyManager().(*dependencyManager)

	suite.NotNil(target.commonDependencyService)
	suite.NotNil(target.sourdoughRecipeDependencyService)
	suite.NotNil(target.sourdoughRecipeScaleDependencyService)
	suite.NotNil(target.flourDependencyService)
}

func TestDependencyManagerTestSuite(t *testing.T) {
	suite.Run(t, new(DependencyManagerTestSuite))
}

type MockValueType struct {
	value string
}

type MockValueIncorrectType struct {
	value int
}

func TestGetFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		key         string
		expectedErr error
	}{
		{
			name:        "Nil context",
			ctx:         nil,
			key:         "mockValue",
			expectedErr: errors.New("context is nil"),
		},
		{
			name:        "Value is nil",
			ctx:         context.Background(),
			key:         "mockValue",
			expectedErr: errors.New("mockValue is nil"),
		},
		{
			name:        "Value is incorrect type",
			ctx:         context.WithValue(context.Background(), "mockValue", MockValueIncorrectType{value: 6}),
			key:         "mockValue",
			expectedErr: errors.New("mockValue is not valid type"),
		},
		{
			name:        "Happy Path",
			ctx:         context.WithValue(context.Background(), "mockValue", MockValueType{value: "a value"}),
			key:         "mockValue",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getFromContext[MockValueType](tt.ctx, tt.key)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsNil(t *testing.T) {
	var nonemptyString = "not nil"
	var nonemptyInt = 3

	testCases := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Nil String", nil, true},
		{"Non Nil String", &nonemptyString, false},
		{"Nil Integer", nil, true},
		{"Non Nil Integer", &nonemptyInt, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isNil(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
