package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeScaleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeScaleServiceTestSuite))
}

type SourdoughRecipeScaleServiceTestSuite struct {
	test.GoMockTestSuite

	ctx                         context.Context
	sourdoughRecipeScaleService *mocks.MockSourdoughRecipeService

	target domain.SourdoughRecipeScaleService
}

func (suite *SourdoughRecipeScaleServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.ctx = context.Background()
	suite.sourdoughRecipeScaleService = mocks.NewMockSourdoughRecipeService(suite.MockCtrl)

	suite.target = test.Must(func() (domain.SourdoughRecipeScaleService, error) {
		return NewSourdoughRecipeScaleService(suite.sourdoughRecipeScaleService)
	})
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScale() {
	dto := createValidDTO(domain.SourdoughRecipeDto{

		RecipeDto: domain.RecipeDto{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
	})

	suite.sourdoughRecipeScaleService.EXPECT().FindById(suite.ctx, dto.Id).
		Return(dto, nil)

	scaledDto, err := suite.target.Scale(suite.ctx, dto.Id, domain.SourdoughRecipeScaleRequestDto{FinalDoughWeight: 985})

	suite.NoError(err)
	suite.Equal(domain.SourdoughRecipeDto{
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
					Amount: 450,
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
			CreatedAt: dto.CreatedAt,
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
					Amount: 23,
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
					Amount: 23,
				},
			},
			Starter: domain.BakerAmountDto{
				Amount: 10,
			},
			Water: domain.BakerAmountDto{
				Amount: 45,
			},
			Amount: domain.BakerAmountDto{
				Amount:          100,
				BakerPercentage: 20,
			},
		},
	}, scaledDto)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScale_WithScaleExists() {
	dto := createValidDTO(domain.SourdoughRecipeDto{
		RecipeDto: domain.RecipeDto{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
	})

	suite.sourdoughRecipeScaleService.EXPECT().FindById(suite.ctx, dto.Id).
		Return(dto, nil)

	scaledDto, err := suite.target.Scale(suite.ctx, dto.Id, domain.SourdoughRecipeScaleRequestDto{FinalDoughWeight: 985})

	suite.NoError(err)
	expectedDto := domain.SourdoughRecipeDto{
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
					Amount: 450,
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
			CreatedAt: dto.CreatedAt,
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
					Amount: 23,
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
					Amount: 23,
				},
			},
			Starter: domain.BakerAmountDto{
				Amount: 10,
			},
			Water: domain.BakerAmountDto{
				Amount: 45,
			},
			Amount: domain.BakerAmountDto{
				Amount:          100,
				BakerPercentage: 20,
			},
		},
	}
	suite.Equal(expectedDto, scaledDto)

	newScaledDto, err := suite.target.Scale(suite.ctx, dto.Id, domain.SourdoughRecipeScaleRequestDto{FinalDoughWeight: 985})

	suite.NoError(err)
	suite.Equal(expectedDto, newScaledDto)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScale_WithErrorOnFind() {
	suite.sourdoughRecipeScaleService.EXPECT().FindById(suite.ctx, gomock.Any()).
		Return(domain.SourdoughRecipeDto{}, assert.AnError)

	scaledDto, err := suite.target.Scale(suite.ctx, test.FirstId, domain.SourdoughRecipeScaleRequestDto{FinalDoughWeight: 985})

	suite.Equal(assert.AnError, err)
	suite.Empty(scaledDto)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleLevain() {
	service := suite.target.(*sourdoughRecipeScaleService)

	flour1 := domain.FlourDto{
		Id:             uuid.New(),
		FlourType:      "flour 1",
		Name:           "flour 1",
		Description:    "flour 1",
		NutritionFacts: domain.NutritionFactsDto{},
	}
	flour2 := domain.FlourDto{
		Id:             uuid.New(),
		FlourType:      "flour 2",
		Name:           "flour 2",
		Description:    "flour 2",
		NutritionFacts: domain.NutritionFactsDto{},
	}
	scaledLevain := service.scaleLevain(1500, domain.SourdoughLevainAgentDto{
		Starter: domain.BakerAmountDto{
			Amount:          20,
			BakerPercentage: 2,
			Name:            "Starter",
		},
		Flour: []domain.FlourAmountDto{
			{
				FlourDto: flour1,
				Amount:   700,
			},
			{
				FlourDto: flour2,
				Amount:   300,
			},
		},
		Water: domain.BakerAmountDto{
			Amount:          500,
			BakerPercentage: 50,
			Name:            "water",
		},
	}, 1000)

	suite.Equal(domain.SourdoughLevainAgentDto{
		Starter: domain.BakerAmountDto{
			Amount:          13,
			BakerPercentage: 2,
			Name:            "Starter",
		},
		Flour: []domain.FlourAmountDto{
			{
				FlourDto: flour1,
				Amount:   467,
			},
			{
				FlourDto: flour2,
				Amount:   200,
			},
		},
		Water: domain.BakerAmountDto{
			Amount:          333,
			BakerPercentage: 50,
			Name:            "water",
		},
	}, scaledLevain)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleRecipeDetails() {
	service := suite.target.(*sourdoughRecipeScaleService)

	scaledRecipeDetails := service.scaleRecipeDetails(2000, domain.RecipeDetailsDto{
		Flour: domain.BakerAmountDto{
			Amount:          1000,
			BakerPercentage: 100,
			Name:            "flour",
		},
		Water: domain.BakerAmountDto{
			Amount:          800,
			BakerPercentage: 80,
			Name:            "Water",
		},
		Levain: domain.BakerAmountDto{
			Amount:          180,
			BakerPercentage: 18,
			Name:            "Levain",
		},
		AdditionalIngredients: domain.BakerAmountDto{
			Amount:          20,
			BakerPercentage: 2,
			Name:            "Salt",
		},
		TotalWeight: 2000,
	}, 1000)

	suite.Equal(domain.RecipeDetailsDto{
		Flour: domain.BakerAmountDto{
			Amount:          500,
			BakerPercentage: 100,
			Name:            "flour",
		},
		Water: domain.BakerAmountDto{
			Amount:          400,
			BakerPercentage: 80,
			Name:            "Water",
		},
		Levain: domain.BakerAmountDto{
			Amount:          90,
			BakerPercentage: 18,
			Name:            "Levain",
		},
		AdditionalIngredients: domain.BakerAmountDto{
			Amount:          10,
			BakerPercentage: 2,
			Name:            "Salt",
		},
		TotalWeight: 1000,
	}, scaledRecipeDetails)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleFlourAmounts() {
	service := suite.target.(*sourdoughRecipeScaleService)
	flourEntity1 := domain.FlourDto{Id: uuid.New()}
	flourEntity2 := domain.FlourDto{Id: uuid.New()}

	amount := service.scaleFlourAmounts(2000, []domain.FlourAmountDto{
		{
			FlourDto: flourEntity1,
			Amount:   1000,
		},
		{
			FlourDto: flourEntity2,
			Amount:   250,
		},
	}, 1000)

	suite.Equal([]domain.FlourAmountDto{
		{
			FlourDto: flourEntity1,
			Amount:   500,
		},
		{
			FlourDto: flourEntity2,
			Amount:   125,
		},
	}, amount)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleFlourAmount() {
	service := suite.target.(*sourdoughRecipeScaleService)
	flourEntity := domain.FlourDto{Id: uuid.New()}

	amount := service.scaleFlourAmount(2000, domain.FlourAmountDto{
		FlourDto: flourEntity,
		Amount:   1000,
	}, 1000)

	suite.Equal(domain.FlourAmountDto{
		FlourDto: flourEntity,
		Amount:   500,
	}, amount)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleBakerAmounts() {
	service := suite.target.(*sourdoughRecipeScaleService)

	amount := service.scaleBakerAmounts(2000, []domain.BakerAmountDto{
		{
			Amount:          1000,
			BakerPercentage: 50,
			Name:            "test",
		},
		{
			Amount:          250,
			BakerPercentage: 25,
			Name:            "test 2",
		},
	}, 1000)

	suite.Equal([]domain.BakerAmountDto{
		{
			Amount:          500,
			BakerPercentage: 50,
			Name:            "test",
		},
		{
			Amount:          125,
			BakerPercentage: 25,
			Name:            "test 2",
		},
	}, amount)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleBakerAmount() {
	service := suite.target.(*sourdoughRecipeScaleService)

	amount := service.scaleBakerAmount(2000, domain.BakerAmountDto{
		Amount:          1000,
		BakerPercentage: 50,
		Name:            "test",
	}, 1000)

	suite.Equal(domain.BakerAmountDto{
		Amount:          500,
		BakerPercentage: 50,
		Name:            "test",
	}, amount)
}

func (suite *SourdoughRecipeScaleServiceTestSuite) TestScaleAmount() {
	service := suite.target.(*sourdoughRecipeScaleService)

	amount := service.scaleAmount(2000, 1000, 1000)

	suite.Equal(float64(500), amount)
}

func TestNewSourdoughRecipeScaleService_WithNilRepository(t *testing.T) {
	service, err := NewSourdoughRecipeScaleService(nil)

	assert.Nil(t, service)
	assert.ErrorContains(t, err, "sourdoughRecipeService cannot be nil")
}
