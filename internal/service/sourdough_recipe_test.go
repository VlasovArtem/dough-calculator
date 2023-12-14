package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	internalErrors "dough-calculator/internal/errors"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeServiceTestSuite))
}

type SourdoughRecipeServiceTestSuite struct {
	test.GoMockTestSuite

	ctx        context.Context
	repository *mocks.MockSourdoughRecipeRepository

	target domain.SourdoughRecipeService
}

func (suite *SourdoughRecipeServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.ctx = context.Background()
	suite.repository = mocks.NewMockSourdoughRecipeRepository(suite.MockCtrl)

	suite.target = test.Must(func() (domain.SourdoughRecipeService, error) {
		return NewSourdoughRecipeService(suite.repository)
	})
}

func (suite *SourdoughRecipeServiceTestSuite) TestCreate() {
	createRequest := generateCreateRequest()

	suite.repository.EXPECT().
		Create(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, entity domain.SourdoughRecipeEntity) (domain.SourdoughRecipeEntity, error) {
			return entity, nil
		})

	dto, err := suite.target.Create(suite.ctx, createRequest)

	suite.NoError(err)
	suite.Equal(createValidDTO(dto), dto)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCreate_WithErrorFromRepository() {
	createRequest := generateCreateRequest()

	suite.repository.EXPECT().
		Create(suite.ctx, gomock.Any()).
		Return(domain.SourdoughRecipeEntity{}, assert.AnError)

	dto, err := suite.target.Create(suite.ctx, createRequest)

	suite.Error(err)
	suite.Equal(err, internalErrors.NewInternalServerErrorWrap(assert.AnError, "failed to create recipe"))
	suite.Empty(dto)
}

func (suite *SourdoughRecipeServiceTestSuite) TestFindById() {
	entity := domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id: uuid.New(),
		},
	}

	suite.repository.EXPECT().
		GetById(suite.ctx, gomock.Any()).
		Return(entity, nil)

	result, err := suite.target.FindById(suite.ctx, entity.Id)

	suite.NoError(err)
	suite.Equal(entity.ToDto(), result)
}

func (suite *SourdoughRecipeServiceTestSuite) TestFindById_WithError() {
	entity := domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id: uuid.New(),
		},
	}

	tests := []struct {
		name                string
		errorFromRepository error
		expectedError       error
	}{
		{
			name:                "with basic error",
			errorFromRepository: assert.AnError,
			expectedError:       internalErrors.NewInternalServerErrorWrap(assert.AnError, "failed to find recipe by id"),
		},
		{
			name:                "with document not found error",
			errorFromRepository: mongo.ErrNoDocuments,
			expectedError:       internalErrors.SourdoughRecipeNotFound(fmt.Sprintf("recipe with id %s not found", entity.Id.String())),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.repository.EXPECT().
				GetById(suite.ctx, entity.Id).
				Return(entity, tt.errorFromRepository)

			_, err := suite.target.FindById(suite.ctx, entity.Id)

			suite.Equal(tt.expectedError, err)
		})
	}
}

func (suite *SourdoughRecipeServiceTestSuite) TestFind() {
	entity := domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id: uuid.New(),
		},
	}

	suite.repository.EXPECT().
		Find(suite.ctx, 10, 0).
		Return([]domain.SourdoughRecipeEntity{entity}, nil)

	result, err := suite.target.Find(suite.ctx, 10, 0)

	suite.NoError(err)
	suite.Equal([]domain.SourdoughRecipeDto{
		{
			RecipeDto: domain.RecipeDto{
				Id:                    entity.Id,
				Water:                 make([]domain.BakerAmountDto, 0),
				Flour:                 make([]domain.FlourAmountDto, 0),
				AdditionalIngredients: make([]domain.BakerAmountDto, 0),
				NutritionFacts:        make(map[string]domain.NutritionFactsDto),
			},
			Levain: domain.SourdoughLevainAgentDto{
				Flour: make([]domain.FlourAmountDto, 0),
			},
		},
	}, result)
}

func (suite *SourdoughRecipeServiceTestSuite) TestFind_WithError() {
	suite.repository.EXPECT().
		Find(suite.ctx, gomock.Any(), gomock.Any()).
		Return([]domain.SourdoughRecipeEntity{}, assert.AnError)

	result, err := suite.target.Find(suite.ctx, 10, 0)

	suite.Empty(result)
	suite.Equal(internalErrors.NewInternalServerErrorWrap(assert.AnError, "failed to find recipes"), err)
}

func (suite *SourdoughRecipeServiceTestSuite) TestSearchByName() {
	entity := domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id: uuid.New(),
		},
	}

	suite.repository.EXPECT().
		SearchByName(suite.ctx, "name").
		Return([]domain.SourdoughRecipeEntity{entity}, nil)

	result, err := suite.target.SearchByName(suite.ctx, "name")

	suite.NoError(err)
	suite.Equal([]domain.SourdoughRecipeDto{entity.ToDto()}, result)
}

func (suite *SourdoughRecipeServiceTestSuite) TestSearchByName_WithError() {
	suite.repository.EXPECT().
		SearchByName(suite.ctx, "name").
		Return([]domain.SourdoughRecipeEntity{}, assert.AnError)

	result, err := suite.target.SearchByName(suite.ctx, "name")

	suite.Empty(result)
	suite.Equal(internalErrors.NewInternalServerErrorWrap(assert.AnError, "failed to search recipes by name"), err)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateRecipeDetails() {
	service := suite.target.(*sourdoughRecipeService)

	request := generateCreateRequest()

	recipeDetails := service.calculateRecipeDetails(request)

	suite.Equal(domain.RecipeDetailsDto{
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
	}, recipeDetails)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateFlourAmount() {
	service := suite.target.(*sourdoughRecipeService)

	amount := service.calculateFlourAmount([]domain.FlourAmountDto{
		{Amount: 900},
		{Amount: 100},
	})

	suite.Equal(domain.BakerAmountDto{
		Amount:          1000,
		BakerPercentage: 100,
	}, amount)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateFlourAmount_WithEmptyIngredients_ShouldReturnEmpty() {
	service := suite.target.(*sourdoughRecipeService)

	tests := []struct {
		name        string
		ingredients []domain.FlourAmountDto
	}{
		{
			name:        "with empty water",
			ingredients: []domain.FlourAmountDto{},
		},
		{
			name:        "with nil water",
			ingredients: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			amount := service.calculateFlourAmount(tt.ingredients)

			suite.Empty(amount)
		})
	}
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateWaterAmount() {
	service := suite.target.(*sourdoughRecipeService)

	amount := service.calculateWaterAmount(domain.BakerAmountDto{Amount: 1000}, []domain.BakerAmountDto{
		{Amount: 700},
		{Amount: 50},
	})

	suite.Equal(domain.BakerAmountDto{
		Amount:          750,
		BakerPercentage: 75,
	}, amount)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateWaterAmount_WithEmptyIngredients_ShouldReturnEmpty() {
	service := suite.target.(*sourdoughRecipeService)

	tests := []struct {
		name        string
		ingredients []domain.BakerAmountDto
	}{
		{
			name:        "with empty water",
			ingredients: []domain.BakerAmountDto{},
		},
		{
			name:        "with nil water",
			ingredients: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			amount := service.calculateWaterAmount(domain.BakerAmountDto{Amount: 1000}, tt.ingredients)

			suite.Empty(amount)
		})
	}
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateAdditionalIngredientsAmount() {
	service := suite.target.(*sourdoughRecipeService)

	amount := service.calculateAdditionalIngredientsAmount(domain.BakerAmountDto{Amount: 1000}, []domain.BakerAmountDto{
		{Amount: 10},
		{Amount: 50},
	})

	suite.Equal(domain.BakerAmountDto{
		Amount:          60,
		BakerPercentage: 6,
	}, amount)
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateAdditionalIngredientsAmount_WithEmptyAdditionalIngredients_ShouldReturnEmpty() {
	service := suite.target.(*sourdoughRecipeService)

	tests := []struct {
		name        string
		ingredients []domain.BakerAmountDto
	}{
		{
			name:        "with empty additional ingredients",
			ingredients: []domain.BakerAmountDto{},
		},
		{
			name:        "with nil additional ingredients",
			ingredients: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			amount := service.calculateAdditionalIngredientsAmount(domain.BakerAmountDto{Amount: 1000}, tt.ingredients)

			suite.Empty(amount)
		})
	}
}

func (suite *SourdoughRecipeServiceTestSuite) TestCalculateTotalWeight() {
	service := suite.target.(*sourdoughRecipeService)

	totalWeight := service.calculateTotalWeight(domain.BakerAmountDto{Amount: 100}, domain.BakerAmountDto{Amount: 13.33}, domain.BakerAmountDto{Amount: 99.99}, domain.BakerAmountDto{Amount: 150})

	suite.Equal(363, totalWeight)
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

func createValidDTO(dto domain.SourdoughRecipeDto) domain.SourdoughRecipeDto {
	return domain.SourdoughRecipeDto{
		RecipeDto: domain.RecipeDto{
			Id:          dto.Id,
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
			CreatedAt: dto.CreatedAt,
			UpdatedAt: nil,
			Yield: domain.RecipeYieldDto{
				Unit:   "loaf",
				Amount: 2,
			},
		},
		Levain: domain.SourdoughLevainAgentDto{
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
			Amount: domain.BakerAmountDto{
				Amount:          200,
				BakerPercentage: 20,
			},
		},
	}
}

func TestNewSourdoughRecipeService_WithNilRepository(t *testing.T) {
	service, err := NewSourdoughRecipeService(nil)

	assert.Nil(t, service)
	assert.ErrorContains(t, err, "repository cannot be nil")
}
