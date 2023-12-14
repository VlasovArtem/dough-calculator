package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

func TestNewApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	configManager := mocks.NewMockConfigManager(mockCtrl)
	configManager.EXPECT().GetConfig().Return(config.Config{
		Application: config.Application{
			Name: "dough-calculator",
		},
	})

	commonDependencyService := mocks.NewMockCommonDependencyService(mockCtrl)
	commonDependencyService.EXPECT().ConfigManager().Return(configManager)

	dependencyMgr := mocks.NewMockDependencyManager(mockCtrl)
	dependencyMgr.EXPECT().Common().Return(commonDependencyService)
	server := &http.Server{}

	var app domain.Application = &application{dependencyManager: dependencyMgr, server: server}

	assert.Equal(t, config.Config{
		Application: config.Application{
			Name: "dough-calculator",
		},
	}, app.Config())
	assert.Equal(t, server, app.Server())
}

func TestApplicationInitializerTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationInitializerTestSuite))
}

type ApplicationInitializerTestSuite struct {
	test.GoMockTestSuite

	configManager                         *mocks.MockConfigManager
	dependencyManager                     *mocks.MockDependencyManager
	commonDependencyService               *mocks.MockCommonDependencyService
	sourdoughRecipeDependencyService      *mocks.MockSourdoughRecipeDependencyService
	sourdoughRecipeScaleDependencyService *mocks.MockSourdoughRecipeScaleDependencyService
	flourDependencyService                *mocks.MockFlourDependencyService

	actuatorHandler             *mocks.MockActuatorHandler
	sourdoughRecipeHandler      *mocks.MockSourdoughRecipeHandler
	sourdoughRecipeScaleHandler *mocks.MockSourdoughRecipeScaleHandler
	flourHandler                *mocks.MockFlourHandler

	target *applicationInitializer
}

func (suite *ApplicationInitializerTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.configManager = mocks.NewMockConfigManager(suite.MockCtrl)
	suite.dependencyManager = mocks.NewMockDependencyManager(suite.MockCtrl)
	suite.commonDependencyService = mocks.NewMockCommonDependencyService(suite.MockCtrl)
	suite.sourdoughRecipeDependencyService = mocks.NewMockSourdoughRecipeDependencyService(suite.MockCtrl)
	suite.sourdoughRecipeScaleDependencyService = mocks.NewMockSourdoughRecipeScaleDependencyService(suite.MockCtrl)
	suite.flourDependencyService = mocks.NewMockFlourDependencyService(suite.MockCtrl)

	suite.actuatorHandler = mocks.NewMockActuatorHandler(suite.MockCtrl)
	suite.sourdoughRecipeHandler = mocks.NewMockSourdoughRecipeHandler(suite.MockCtrl)
	suite.sourdoughRecipeScaleHandler = mocks.NewMockSourdoughRecipeScaleHandler(suite.MockCtrl)
	suite.flourHandler = mocks.NewMockFlourHandler(suite.MockCtrl)

	suite.target = &applicationInitializer{dependencyManager: suite.dependencyManager}
}

func (suite *ApplicationInitializerTestSuite) TestInitialize() {
	suite.dependencyManager.EXPECT().Initialize(gomock.Any()).Return(nil)
	suite.dependencyManager.EXPECT().Common().Return(suite.commonDependencyService).AnyTimes()
	suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager).AnyTimes()
	suite.configManager.EXPECT().GetConfig().Return(config.Config{
		Application: config.Application{
			Rest: config.Rest{
				ContextPath: "/api",
			},
		},
	}).AnyTimes()

	suite.commonDependencyService.EXPECT().Actuator().Return(suite.actuatorHandler)
	suite.actuatorHandler.EXPECT().Health().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	suite.dependencyManager.EXPECT().SourdoughRecipe().Return(suite.sourdoughRecipeDependencyService)
	suite.sourdoughRecipeDependencyService.EXPECT().Router().Return(suite.sourdoughRecipeHandler)
	suite.sourdoughRecipeHandler.EXPECT().Create().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.sourdoughRecipeHandler.EXPECT().Find().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.sourdoughRecipeHandler.EXPECT().FindById().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.sourdoughRecipeHandler.EXPECT().Search().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	suite.dependencyManager.EXPECT().SourdoughRecipeScale().Return(suite.sourdoughRecipeScaleDependencyService)
	suite.sourdoughRecipeScaleDependencyService.EXPECT().Router().Return(suite.sourdoughRecipeScaleHandler)
	suite.sourdoughRecipeScaleHandler.EXPECT().Scale().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	suite.dependencyManager.EXPECT().Flour().Return(suite.flourDependencyService)
	suite.flourDependencyService.EXPECT().Router().Return(suite.flourHandler)
	suite.flourHandler.EXPECT().Create().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.flourHandler.EXPECT().Find().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.flourHandler.EXPECT().FindById().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	suite.flourHandler.EXPECT().Search().
		Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	app, err := suite.target.Initialize()

	assert.NotNil(suite.T(), app)
	assert.NoError(suite.T(), err)
}

func (suite *ApplicationInitializerTestSuite) TestInitialize_WithError() {
	tests := []struct {
		name                 string
		mocksConfigured      func()
		expectedErrorMessage string
	}{
		{
			name: "failed to initialize dependency",
			mocksConfigured: func() {
				suite.dependencyManager.EXPECT().Initialize(gomock.Any()).Return(assert.AnError)
			},
			expectedErrorMessage: "failed to initialize dependency",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mocksConfigured()

			app, err := suite.target.Initialize()

			assert.Nil(suite.T(), app)
			assert.ErrorContains(suite.T(), err, tt.expectedErrorMessage)
		})
	}
}

func (suite *ApplicationInitializerTestSuite) TestInitializeRouter() {
	suite.dependencyManager.EXPECT().Common().Return(suite.commonDependencyService).AnyTimes()
	suite.commonDependencyService.EXPECT().ConfigManager().Return(suite.configManager).AnyTimes()
	suite.configManager.EXPECT().GetConfig().
		Return(config.Config{
			Application: config.Application{
				Rest: config.Rest{
					ContextPath: "/api",
				},
			},
		})

	defaultHandlerProvider := func(message string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(message))
			if err != nil {
				suite.Require().Error(err)
			}
		}
	}

	suite.commonDependencyService.EXPECT().Actuator().Return(suite.actuatorHandler)
	suite.actuatorHandler.EXPECT().Health().Return(defaultHandlerProvider("health ok"))

	suite.dependencyManager.EXPECT().SourdoughRecipe().Return(suite.sourdoughRecipeDependencyService)
	suite.sourdoughRecipeDependencyService.EXPECT().Router().Return(suite.sourdoughRecipeHandler)
	suite.sourdoughRecipeHandler.EXPECT().Create().
		Return(defaultHandlerProvider("create sourdough recipe ok"))
	suite.sourdoughRecipeHandler.EXPECT().Find().
		Return(defaultHandlerProvider("find sourdough recipe ok"))
	suite.sourdoughRecipeHandler.EXPECT().FindById().
		Return(defaultHandlerProvider("find by id sourdough recipe ok"))
	suite.sourdoughRecipeHandler.EXPECT().Search().
		Return(defaultHandlerProvider("search sourdough recipe ok"))

	suite.dependencyManager.EXPECT().SourdoughRecipeScale().Return(suite.sourdoughRecipeScaleDependencyService)
	suite.sourdoughRecipeScaleDependencyService.EXPECT().Router().Return(suite.sourdoughRecipeScaleHandler)
	suite.sourdoughRecipeScaleHandler.EXPECT().Scale().
		Return(defaultHandlerProvider("scale sourdough recipe ok"))

	suite.dependencyManager.EXPECT().Flour().Return(suite.flourDependencyService)
	suite.flourDependencyService.EXPECT().Router().Return(suite.flourHandler)
	suite.flourHandler.EXPECT().Create().
		Return(defaultHandlerProvider("create flour ok"))
	suite.flourHandler.EXPECT().Find().
		Return(defaultHandlerProvider("find flour ok"))
	suite.flourHandler.EXPECT().FindById().
		Return(defaultHandlerProvider("find by id flour ok"))
	suite.flourHandler.EXPECT().Search().
		Return(defaultHandlerProvider("search flour ok"))

	router := suite.target.initializeRouter()

	suite.Run("health", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/actuator/health", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("health ok", resp.Body.String())
	})

	suite.Run("create sourdough recipe", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/recipe/sourdough", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("create sourdough recipe ok", resp.Body.String())
	})

	suite.Run("find sourdough recipe", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/recipe/sourdough", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("find sourdough recipe ok", resp.Body.String())
	})

	suite.Run("find by id sourdough recipe", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/recipe/sourdough/1", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("find by id sourdough recipe ok", resp.Body.String())
	})

	suite.Run("search sourdough recipe", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/recipe/sourdough/search", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("search sourdough recipe ok", resp.Body.String())
	})

	suite.Run("scale sourdough recipe", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/recipe/sourdough/1/scale", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("scale sourdough recipe ok", resp.Body.String())
	})

	suite.Run("create flour", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/flour", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("create flour ok", resp.Body.String())
	})

	suite.Run("find flour", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/flour", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("find flour ok", resp.Body.String())
	})

	suite.Run("find by id flour", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/flour/1", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("find by id flour ok", resp.Body.String())
	})

	suite.Run("search flour", func() {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/flour/search", nil))

		suite.Equal(http.StatusOK, resp.Code)
		suite.Equal("search flour ok", resp.Body.String())
	})
}

func (suite *ApplicationInitializerTestSuite) TestApplicationInitializer_WithError() {
	tests := []struct {
		name                 string
		dependencyManager    func() domain.DependencyManager
		expectedErrorMessage string
	}{
		{
			name: "failed to initialize dependency manager",
			dependencyManager: func() domain.DependencyManager {
				suite.dependencyManager.EXPECT().Initialize(gomock.Any()).Return(assert.AnError)
				return suite.dependencyManager
			},
			expectedErrorMessage: "failed to initialize dependency",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			appInitializer := &applicationInitializer{dependencyManager: tt.dependencyManager()}

			initialize, err := appInitializer.Initialize()
			suite.ErrorContains(err, tt.expectedErrorMessage)
			suite.Nil(initialize)
		})
	}
}
