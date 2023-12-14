package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	internalErrors "dough-calculator/internal/errors"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeHandlerTestSuite))
}

type SourdoughRecipeHandlerTestSuite struct {
	test.GoMockTestSuite

	service *mocks.MockSourdoughRecipeService

	target domain.SourdoughRecipeHandler
}

func (suite *SourdoughRecipeHandlerTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.service = mocks.NewMockSourdoughRecipeService(suite.MockCtrl)

	suite.target = test.Must(func() (domain.SourdoughRecipeHandler, error) {
		return NewSourdoughRecipeHandler(suite.service)
	})
}

func (suite *SourdoughRecipeHandlerTestSuite) TestCreate() {
	request := generateCreateRequest()

	suite.service.EXPECT().
		Create(gomock.Any(), request).
		Return(createSourdoughRecipe(), nil)

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/", suite.target.Create())

	req, err := http.NewRequest("POST", "/", buffer)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFileAsObject[domain.SourdoughRecipeDto](suite.T(), resp, http.StatusCreated, "testdata/sourdough_recipe_response.json")
}

func (suite *SourdoughRecipeHandlerTestSuite) TestCreate_WithInvalidRequest() {
	req := httptest.NewRequest("POST", "http://testing", bytes.NewBuffer([]byte("invalid body")))
	resp := httptest.NewRecorder()

	suite.target.Create().ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": -1,
			"error_details": "error while decoding request body: invalid character 'i' looking for beginning of value",
			"error_message": "internal server error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusInternalServerError, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestCreate_WithErrorOnCreate() {
	request := generateCreateRequest()

	suite.service.EXPECT().
		Create(gomock.Any(), request).
		Return(domain.SourdoughRecipeDto{}, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/", suite.target.Create())

	req, err := http.NewRequest("POST", "/", buffer)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 123,
			"error_details": "error 'test'",
			"error_message": "error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFindById() {
	recipe := createSourdoughRecipe()

	suite.service.EXPECT().FindById(gomock.Any(), recipe.Id).
		Return(recipe, nil)

	router := chi.NewRouter()
	router.
		Get("/recipe/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", fmt.Sprintf("/recipe/%s", recipe.Id), nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/sourdough_recipe_response.json")
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFindById_WithoutParam() {
	router := chi.NewRouter()
	router.
		Get("/recipe", suite.target.FindById())

	req, err := http.NewRequest("GET", "/recipe", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 10001,
			"error_details": "id is required",
			"error_message": "id is required"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFindById_WithInvalidIdParam() {
	router := chi.NewRouter()
	router.
		Get("/recipe/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", "/recipe/invalid", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 10002,
			"error_details": "id is not valid",
			"error_message": "id is not valid"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFindById_WithErrorOnFindById() {
	recipe := createSourdoughRecipe()

	suite.service.EXPECT().FindById(gomock.Any(), recipe.Id).
		Return(domain.SourdoughRecipeDto{}, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	router := chi.NewRouter()
	router.
		Get("/recipe/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", fmt.Sprintf("/recipe/%s", recipe.Id), nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 123,
			"error_details": "error 'test'",
			"error_message": "error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFind() {
	recipes := []domain.SourdoughRecipeDto{createSourdoughRecipe()}

	suite.service.EXPECT().Find(gomock.Any(), 1, 10).Return(recipes, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(PageInput{})).
		Get("/find", suite.target.Find())

	req, err := http.NewRequest("GET", "/find?offset=1&limit=10", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/sourdough_recipes_response.json")
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFind_WithDefaultParameters() {
	recipes := []domain.SourdoughRecipeDto{createSourdoughRecipe()}

	suite.service.EXPECT().Find(gomock.Any(), 0, 25).Return(recipes, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(PageInput{})).
		Get("/", suite.target.Find())

	req, err := http.NewRequest("GET", "/", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/sourdough_recipes_response.json")
}

func (suite *SourdoughRecipeHandlerTestSuite) TestFind_WithErrorOnFind() {
	suite.service.EXPECT().Find(gomock.Any(), 0, 25).
		Return(nil, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(PageInput{})).
		Get("/", suite.target.Find())

	req, err := http.NewRequest("GET", "/", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 123,
			"error_details": "error 'test'",
			"error_message": "error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *SourdoughRecipeHandlerTestSuite) TestSearch() {
	recipes := []domain.SourdoughRecipeDto{createSourdoughRecipe()}

	suite.service.EXPECT().SearchByName(gomock.Any(), "test name").
		Return(recipes, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(SearchRecipeInput{})).
		Get("/search", suite.target.Search())

	req, err := http.NewRequest("GET", "/search?name=test%20name", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/sourdough_recipes_response.json")
}

func (suite *SourdoughRecipeHandlerTestSuite) TestSearch_WithErrorOnSearch() {
	suite.service.EXPECT().SearchByName(gomock.Any(), "test name").
		Return(nil, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(SearchRecipeInput{})).
		Get("/search", suite.target.Search())

	req, err := http.NewRequest("GET", "/search?name=test%20name", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": 123,
			"error_details": "error 'test'",
			"error_message": "error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func TestNewSourdoughRecipeHandler_WithNilService(t *testing.T) {
	handler, err := NewSourdoughRecipeHandler(nil)

	assert.ErrorContains(t, err, "service cannot be nil")
	assert.Nil(t, handler)
}

func generateCreateRequest() domain.CreateSourdoughRecipeRequest {
	return domain.CreateSourdoughRecipeRequest{
		Name:        "test recipe",
		Description: "test recipe description",
		Flour: []domain.FlourAmountDto{
			{
				FlourDto: domain.FlourDto{
					Id:          test.FirstId,
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
					Id:          test.SecondId,
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
		Levain: domain.SourdoughLevainAgentDto{
			Amount: domain.BakerAmountDto{
				Amount:          200,
				BakerPercentage: 20,
			},
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          test.FirstId,
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
						Id:          test.SecondId,
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
			Starter: domain.BakerAmountDto{
				Amount: 20,
			},
			Water: domain.BakerAmountDto{
				Amount: 90,
			},
		},
		AdditionalIngredients: []domain.BakerAmountDto{
			{
				Amount:          20,
				BakerPercentage: 2,
				Name:            "Salt",
			},
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
			Unit:   "loaf",
			Amount: 2,
		},
	}
}

func createSourdoughRecipe() domain.SourdoughRecipeDto {
	return domain.SourdoughRecipeDto{
		RecipeDto: domain.RecipeDto{
			Id:          test.ThirdId,
			Name:        "test recipe",
			Description: "test recipe description",
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          test.FirstId,
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
						Id:          test.SecondId,
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
			CreatedAt: test.Date,
			UpdatedAt: nil,
			Yield: domain.RecipeYieldDto{
				Unit:   "loaf",
				Amount: 2,
			},
		},
		Levain: domain.SourdoughLevainAgentDto{
			Amount: domain.BakerAmountDto{
				Amount:          200,
				BakerPercentage: 20,
			},
			Flour: []domain.FlourAmountDto{
				{
					FlourDto: domain.FlourDto{
						Id:          test.FirstId,
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
						Id:          test.SecondId,
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
			Starter: domain.BakerAmountDto{
				Amount: 20,
			},
			Water: domain.BakerAmountDto{
				Amount: 90,
			},
		},
	}
}
