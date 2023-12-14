//go:build integration && docker

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"dough-calculator/internal/app"
	"dough-calculator/internal/config"
	"dough-calculator/internal/domain"
	"dough-calculator/internal/test"
)

func TestApplicationTestSuite(t *testing.T) {
	suite.Run(t, &ApplicationTestSuite{
		MongoDBServiceDockerIntegrationTestSuite: test.NewMongoDBServiceDockerIntegrationTestSuite(dockerStarter),
	})
}

type ApplicationTestSuite struct {
	test.MongoDBServiceDockerIntegrationTestSuite

	config config.Config

	initializer domain.ApplicationInitializer
	application domain.Application
	client      *Client
}

func (suite *ApplicationTestSuite) SetupSuite() {
	suite.MongoDBServiceDockerIntegrationTestSuite.SetupSuite()

	configPath := suite.createConfigFile()

	err := os.Setenv("CONFIG_PATH", configPath)
	suite.Require().NoError(err)

	suite.initializer = app.NewApplicationInitializer()

	suite.application = test.Must(func() (domain.Application, error) {
		return suite.initializer.Initialize()
	})

	server := suite.application.Server()

	listener := test.Must(func() (net.Listener, error) {
		return net.Listen("tcp", ":0")
	})

	go func() {
		fmt.Println("Server is ready to handle requests at", listener.Addr())
		if err := server.Serve(listener); err != nil {
			suite.Require().NoError(err)
		}
	}()

	addr := listener.Addr().(*net.TCPAddr)

	suite.isListenerReady(listener)

	suite.client = test.Must(func() (client *Client, err error) {
		return NewClient(fmt.Sprintf("http://localhost:%d", addr.Port))
	})
}

func (suite *ApplicationTestSuite) TearDownTest() {
	err := suite.Drop("dough-calculator", "sourdough-recipes")
	suite.Require().NoError(err)
}

func (suite *ApplicationTestSuite) isListenerReady(listener net.Listener) {
	var serviceIsReady bool
	for i := 0; i < 10; i++ {
		_, serviceErr := net.Dial("tcp", listener.Addr().String())
		if serviceErr == nil {
			serviceIsReady = true
			break
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
	suite.Require().True(serviceIsReady, "service is not ready")
}

func (suite *ApplicationTestSuite) TestApplication_Health() {
	health, err := suite.client.Health(context.Background())
	suite.Require().NoError(err)

	suite.Equal(http.StatusOK, health.StatusCode)
}

func (suite *ApplicationTestSuite) TestApplication_CreateSourdoughRecipe() {
	requestFile, err := os.OpenFile("testdata/sourdough_recipe_create_request.json", os.O_RDONLY, 0644)
	suite.Require().NoError(err)

	response, err := suite.client.CreateSourdoughRecipeWithBody(context.Background(), "application/json", requestFile)
	suite.Require().NoError(err)

	suite.Equal(http.StatusCreated, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	suite.Require().NoError(err)

	fmt.Println(string(bodyBytes))

	var actualResponse domain.SourdoughRecipeDto
	err = json.Unmarshal(bodyBytes, &actualResponse)
	suite.Require().NoError(err)

	suite.Equal(domain.SourdoughRecipeDto{
		RecipeDto: domain.RecipeDto{
			Id:          actualResponse.Id,
			Name:        "test recipe",
			Description: "test recipe description",
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("4487b1c1-672e-425c-bacb-1deb377f0c65"),
						FlourType:   "test first flour type",
						Name:        "test first flour name",
						Description: "test first flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 1,
							Fat:      1,
							Carbs:    1,
							Protein:  1,
							Fiber:    1,
						},
					},
					Amount: 900,
				},
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("1126e515-b2e9-47e5-990d-ad3d8c0f7c98"),
						FlourType:   "test second flour type",
						Name:        "test second flour name",
						Description: "test second flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 2,
							Fat:      2,
							Carbs:    2,
							Protein:  2,
							Fiber:    2,
						},
					},
					Amount: 100,
				},
			},
			Water: []domain.BakerAmountDto{
				{
					Amount:          700,
					BakerPercentage: 70,
					Name:            "Water 1",
				},
				{
					Amount:          50,
					BakerPercentage: 5,
					Name:            "Water 2",
				},
			},
			AdditionalIngredients: []domain.BakerAmountDto{
				{
					Amount:          20,
					BakerPercentage: 2,
					Name:            "Salt",
				},
			},
			Details: domain.RecipeDetailsDto{
				Flour: domain.BakerAmountDto{
					Amount:          1000,
					BakerPercentage: 100,
				},
				Water: domain.BakerAmountDto{
					Amount:          750,
					BakerPercentage: 75,
				},
				Levain: domain.BakerAmountDto{
					Amount:          200,
					BakerPercentage: 20,
				},
				AdditionalIngredients: domain.BakerAmountDto{
					Amount:          20,
					BakerPercentage: 2,
				},
				TotalWeight: 1970,
			},
			NutritionFacts: map[string]domain.NutritionFactsDto{
				"100g": {
					Calories: 1,
					Fat:      1,
					Carbs:    1,
					Protein:  1,
					Fiber:    1,
				},
			},
			Yield: domain.RecipeYieldDto{
				Amount: 2,
				Unit:   "loaf",
			},
			CreatedAt: actualResponse.CreatedAt,
		},
		Levain: domain.SourdoughLevainAgentDto{
			Amount: domain.BakerAmountDto{
				Amount:          200,
				BakerPercentage: 20,
			},
			Starter: domain.BakerAmountDto{
				Amount: 20,
			},
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("4487b1c1-672e-425c-bacb-1deb377f0c65"),
						FlourType:   "test first flour type",
						Name:        "test first flour name",
						Description: "test first flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 1,
							Fat:      1,
							Carbs:    1,
							Protein:  1,
							Fiber:    1,
						},
					},
					Amount: 45,
				},
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("1126e515-b2e9-47e5-990d-ad3d8c0f7c98"),
						FlourType:   "test second flour type",
						Name:        "test second flour name",
						Description: "test second flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 2,
							Fat:      2,
							Carbs:    2,
							Protein:  2,
							Fiber:    2,
						},
					},
					Amount: 45,
				},
			},
			Water: domain.BakerAmountDto{
				Amount: 90,
			},
		},
	}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindSourdoughRecipeById() {
	expectedResponse, err := suite.createSourdoughRecipe()

	sourdoughRecipeResponse, err := suite.client.FindSourdoughRecipeById(context.Background(), expectedResponse.Id)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	var actualResponse domain.SourdoughRecipeDto
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse.CreatedAt

	suite.Equal(expectedResponse, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindSourdoughRecipe() {
	expectedResponse, err := suite.createSourdoughRecipe()
	suite.Require().NoError(err)

	sourdoughRecipeResponse, err := suite.client.FindSourdoughRecipe(context.Background(), &FindSourdoughRecipeParams{
		Offset: intRef(0),
		Limit:  intRef(10),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	var actualResponse []domain.SourdoughRecipeDto
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse[0].CreatedAt

	suite.Equal([]domain.SourdoughRecipeDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindSourdoughRecipe_WithDefaultParameters() {
	expectedResponse, err := suite.createSourdoughRecipe()

	sourdoughRecipeResponse, err := suite.client.FindSourdoughRecipe(context.Background(), nil)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	actualResponse := make([]domain.SourdoughRecipeDto, 1)
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse[0].CreatedAt

	suite.Equal([]domain.SourdoughRecipeDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchSourdoughRecipe() {
	expectedResponse, err := suite.createSourdoughRecipe()

	sourdoughRecipeResponse, err := suite.client.SearchSourdoughRecipe(context.Background(), &SearchSourdoughRecipeParams{
		Name: stringRef(expectedResponse.Name),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	actualResponse := make([]domain.SourdoughRecipeDto, 1)
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse[0].CreatedAt

	suite.Equal([]domain.SourdoughRecipeDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchSourdoughRecipe_WithPattern() {
	expectedResponse, err := suite.createSourdoughRecipe()

	sourdoughRecipeResponse, err := suite.client.SearchSourdoughRecipe(context.Background(), &SearchSourdoughRecipeParams{
		Name: stringRef("test.*"),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	actualResponse := make([]domain.SourdoughRecipeDto, 1)
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse[0].CreatedAt

	suite.Equal([]domain.SourdoughRecipeDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchSourdoughRecipe_WithoutName() {
	expectedResponse, err := suite.createSourdoughRecipe()

	sourdoughRecipeResponse, err := suite.client.SearchSourdoughRecipe(context.Background(), nil)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	actualResponse := make([]domain.SourdoughRecipeDto, 1)
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	expectedResponse.CreatedAt = actualResponse[0].CreatedAt

	suite.Equal([]domain.SourdoughRecipeDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_ScaleSourdoughRecipe() {
	expectedResponse, err := suite.createSourdoughRecipe()
	requestFile, err := os.OpenFile("testdata/sourdough_recipe_scale_request.json", os.O_RDONLY, 0644)
	suite.Require().NoError(err)

	sourdoughRecipeResponse, err := suite.client.ScaleSourdoughRecipeWithBody(context.Background(), expectedResponse.Id, "application/json", requestFile)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, sourdoughRecipeResponse.StatusCode)

	var actualResponse domain.SourdoughRecipeDto
	err = json.NewDecoder(sourdoughRecipeResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal(domain.SourdoughRecipeDto{
		RecipeDto: domain.RecipeDto{
			Id:          actualResponse.Id,
			Name:        "test recipe",
			Description: "test recipe description",
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("4487b1c1-672e-425c-bacb-1deb377f0c65"),
						FlourType:   "test first flour type",
						Name:        "test first flour name",
						Description: "test first flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 1,
							Fat:      1,
							Carbs:    1,
							Protein:  1,
							Fiber:    1,
						},
					},
					Amount: 450,
				},
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("1126e515-b2e9-47e5-990d-ad3d8c0f7c98"),
						FlourType:   "test second flour type",
						Name:        "test second flour name",
						Description: "test second flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 2,
							Fat:      2,
							Carbs:    2,
							Protein:  2,
							Fiber:    2,
						},
					},
					Amount: 50,
				},
			},
			Water: []domain.BakerAmountDto{
				{
					Amount:          350,
					BakerPercentage: 70,
					Name:            "Water 1",
				},
				{
					Amount:          25,
					BakerPercentage: 5,
					Name:            "Water 2",
				},
			},
			AdditionalIngredients: []domain.BakerAmountDto{
				{
					Amount:          10,
					BakerPercentage: 2,
					Name:            "Salt",
				},
			},
			Details: domain.RecipeDetailsDto{
				Flour: domain.BakerAmountDto{
					Amount:          500,
					BakerPercentage: 100,
				},
				Water: domain.BakerAmountDto{
					Amount:          375,
					BakerPercentage: 75,
				},
				Levain: domain.BakerAmountDto{
					Amount:          100,
					BakerPercentage: 20,
				},
				AdditionalIngredients: domain.BakerAmountDto{
					Amount:          10,
					BakerPercentage: 2,
				},
				TotalWeight: 985,
			},
			NutritionFacts: map[string]domain.NutritionFactsDto{
				"100g": {
					Calories: 1,
					Fat:      1,
					Carbs:    1,
					Protein:  1,
					Fiber:    1,
				},
			},
			CreatedAt: actualResponse.CreatedAt,
		},
		Levain: domain.SourdoughLevainAgentDto{
			Amount: domain.BakerAmountDto{
				Amount:          100,
				BakerPercentage: 20,
			},
			Starter: domain.BakerAmountDto{
				Amount: 10,
			},
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("4487b1c1-672e-425c-bacb-1deb377f0c65"),
						FlourType:   "test first flour type",
						Name:        "test first flour name",
						Description: "test first flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 1,
							Fat:      1,
							Carbs:    1,
							Protein:  1,
							Fiber:    1,
						},
					},
					Amount: 23,
				},
				{
					FlourDto: domain.FlourDto{
						Id:          uuid.MustParse("1126e515-b2e9-47e5-990d-ad3d8c0f7c98"),
						FlourType:   "test second flour type",
						Name:        "test second flour name",
						Description: "test second flour description",
						NutritionFacts: domain.NutritionFactsDto{
							Calories: 2,
							Fat:      2,
							Carbs:    2,
							Protein:  2,
							Fiber:    2,
						},
					},
					Amount: 23,
				},
			},
			Water: domain.BakerAmountDto{
				Amount: 45,
			},
		},
	}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_CreateFlour() {
	requestFile, err := os.OpenFile("testdata/flour_create_request.json", os.O_RDONLY, 0644)
	suite.Require().NoError(err)

	response, err := suite.client.CreateFlourWithBody(context.Background(), "application/json", requestFile)
	suite.Require().NoError(err)

	suite.Equal(http.StatusCreated, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	suite.Require().NoError(err)

	fmt.Println(string(bodyBytes))

	var actualResponse domain.FlourDto
	err = json.Unmarshal(bodyBytes, &actualResponse)
	suite.Require().NoError(err)

	suite.Equal(domain.FlourDto{
		Id:          actualResponse.Id,
		FlourType:   "Wheat",
		Name:        "Whole Wheat Flour",
		Description: "Whole wheat flour made from 100% whole wheat grains.",
		NutritionFacts: domain.NutritionFactsDto{
			Calories: 407,
			Fat:      2.2,
			Carbs:    86.4,
			Protein:  16.4,
			Fiber:    12.2,
		},
	}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindFlourById() {
	expectedResponse, err := suite.createFlour()

	clientResponse, err := suite.client.FindFlourById(context.Background(), expectedResponse.Id)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	var actualResponse domain.FlourDto
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal(expectedResponse, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindFlour() {
	expectedResponse, err := suite.createFlour()
	suite.Require().NoError(err)

	clientResponse, err := suite.client.FindFlours(context.Background(), &FindFloursParams{
		Offset: intRef(0),
		Limit:  intRef(10),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	var actualResponse []domain.FlourDto
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal([]domain.FlourDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_FindFlour_WithDefaultParameters() {
	expectedResponse, err := suite.createFlour()

	clientResponse, err := suite.client.FindFlours(context.Background(), nil)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	actualResponse := make([]domain.FlourDto, 1)
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal([]domain.FlourDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchFlour() {
	expectedResponse, err := suite.createFlour()

	clientResponse, err := suite.client.SearchFlour(context.Background(), &SearchFlourParams{
		Name: stringRef(expectedResponse.Name),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	actualResponse := make([]domain.FlourDto, 1)
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal([]domain.FlourDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchFlour_WithPattern() {
	expectedResponse, err := suite.createFlour()

	clientResponse, err := suite.client.SearchFlour(context.Background(), &SearchFlourParams{
		Name: stringRef("Whole.*"),
	})

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	actualResponse := make([]domain.FlourDto, 1)
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal([]domain.FlourDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) TestApplication_SearchFlour_WithoutName() {
	expectedResponse, err := suite.createFlour()

	clientResponse, err := suite.client.SearchFlour(context.Background(), nil)

	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, clientResponse.StatusCode)

	actualResponse := make([]domain.FlourDto, 1)
	err = json.NewDecoder(clientResponse.Body).Decode(&actualResponse)
	suite.Require().NoError(err)

	suite.Equal([]domain.FlourDto{expectedResponse}, actualResponse)
}

func (suite *ApplicationTestSuite) createSourdoughRecipe() (domain.SourdoughRecipeDto, error) {
	requestFile, err := os.OpenFile("testdata/sourdough_recipe_create_request.json", os.O_RDONLY, 0644)
	suite.Require().NoError(err)

	response, err := suite.client.CreateSourdoughRecipeWithBody(context.Background(), "application/json", requestFile)
	suite.Require().NoError(err)

	suite.Equal(http.StatusCreated, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	suite.Require().NoError(err)

	var responseDto domain.SourdoughRecipeDto
	err = json.Unmarshal(bodyBytes, &responseDto)
	suite.Require().NoError(err)

	return responseDto, err
}

func (suite *ApplicationTestSuite) createFlour() (domain.FlourDto, error) {
	requestFile, err := os.OpenFile("testdata/flour_create_request.json", os.O_RDONLY, 0644)
	suite.Require().NoError(err)

	response, err := suite.client.CreateFlourWithBody(context.Background(), "application/json", requestFile)
	suite.Require().NoError(err)

	suite.Equal(http.StatusCreated, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	suite.Require().NoError(err)

	var responseDto domain.FlourDto
	err = json.Unmarshal(bodyBytes, &responseDto)
	suite.Require().NoError(err)

	return responseDto, err
}

func (suite *ApplicationTestSuite) createConfigFile() string {
	tempDir := suite.T().TempDir()

	configFilePath := tempDir + "/config.yaml"

	dbConfig := suite.GetConfig()

	suite.config = config.Config{
		Application: config.Application{
			Name: "loan",
			Rest: config.Rest{
				Server:               ":0",
				ContextPath:          "/v1",
				ReadTimeout:          10,
				WriteTimeout:         10,
				IdleTimeout:          10,
				GraceShutdownTimeout: 10,
			},
		},
		Database: dbConfig,
	}

	cfgBytes, err := yaml.Marshal(suite.config)
	suite.Require().NoError(err)
	err = os.WriteFile(configFilePath, cfgBytes, 0644)
	suite.Require().NoError(err)

	return configFilePath
}

func intRef(v int) *int {
	return &v
}

func stringRef(v string) *string {
	return &v
}
