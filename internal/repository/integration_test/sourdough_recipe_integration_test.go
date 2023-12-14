//go:build integration && docker

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/repository"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &SourdoughRecipeRepositoryTestSuite{
		MongoDBServiceDockerIntegrationTestSuite: test.NewMongoDBServiceDockerIntegrationTestSuite(dockerStarter),
	})
}

type SourdoughRecipeRepositoryTestSuite struct {
	test.MongoDBServiceDockerIntegrationTestSuite

	mockMongoDbService *mocks.MockMongoDBService

	target domain.SourdoughRecipeRepository
}

func (suite *SourdoughRecipeRepositoryTestSuite) SetupSuite() {
	suite.MongoDBServiceDockerIntegrationTestSuite.SetupSuite()

	suite.target = test.Must(func() (domain.SourdoughRecipeRepository, error) {
		return repository.NewSourdoughRecipeRepository(suite.Stub)
	})
}

func (suite *SourdoughRecipeRepositoryTestSuite) AfterTest(suiteName, testName string) {
	err := suite.Drop(repository.SourdoughRecipeDatabase, repository.SourdoughRecipeCollection)
	suite.Require().NoError(err)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestCreate() {
	expected := generateSourdoughRecipeEntity()

	actual, err := suite.target.Create(context.Background(), expected)

	suite.NoError(err)
	suite.Equal(expected, actual)

	var saved domain.SourdoughRecipeEntity

	err = suite.MStub().MustGetCollection(repository.SourdoughRecipeDatabase, repository.SourdoughRecipeCollection).
		FindOne(context.Background(), bson.D{{"_id", expected.Id}}).
		Decode(&saved)

	suite.Equal(expected, saved)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestCreate_WithEntityExists_ShouldReturnError() {
	expected := generateSourdoughRecipeEntity()

	actual, err := suite.target.Create(context.Background(), expected)

	suite.NoError(err)
	suite.Equal(expected, actual)

	_, err = suite.target.Create(context.Background(), expected)

	suite.ErrorContains(err, "failed to insert sourdough recipe")
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestGetById() {
	expected := generateSourdoughRecipeEntity()

	_, err := suite.target.Create(context.Background(), expected)
	suite.NoError(err)

	actual, err := suite.target.GetById(context.Background(), expected.Id)

	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestGetById_WithEntityNotFound_ShouldReturnError() {
	expected := generateSourdoughRecipeEntity()

	_, err := suite.target.GetById(context.Background(), expected.Id)

	suite.ErrorContains(err, "failed to find recipe")
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestFind() {
	first := generateSourdoughRecipeEntity()
	first.CreatedAt = time.Now().Add(-time.Hour).Truncate(time.Second).UTC()
	_, err := suite.target.Create(context.Background(), first)
	suite.Require().NoError(err)

	second := generateSourdoughRecipeEntity()
	second.CreatedAt = time.Now().Truncate(time.Second).UTC()
	_, err = suite.target.Create(context.Background(), second)
	suite.Require().NoError(err)

	actual, err := suite.target.Find(context.Background(), 0, 1)

	suite.NoError(err)
	suite.Equal([]domain.SourdoughRecipeEntity{second}, actual)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestFind_WithEmptyData_ShouldReturnNil() {
	actual, err := suite.target.Find(context.Background(), 1, 0)

	suite.NoError(err)
	suite.Nil(actual)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestFindByName() {
	entity := generateSourdoughRecipeEntity()
	_, err := suite.target.Create(context.Background(), entity)
	suite.Require().NoError(err)

	actual, err := suite.target.SearchByName(context.Background(), entity.Name)

	suite.NoError(err)
	suite.Equal([]domain.SourdoughRecipeEntity{entity}, actual)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestFindByName_WithEntityNotExists_ShouldReturnEmptyEntity() {
	actual, err := suite.target.SearchByName(context.Background(), "missing")

	suite.NoError(err)
	suite.Nil(actual)
}

func generateSourdoughRecipeEntity() domain.SourdoughRecipeEntity {
	id := uuid.New()

	flour := domain.FlourEntity{
		Id:          uuid.New(),
		FlourType:   "test",
		Name:        "test-flour",
		Description: "test flour description",
		NutritionFacts: domain.NutritionFacts{
			Calories: 100,
			Fat:      1,
			Carbs:    1,
			Protein:  2.5,
			Fiber:    1,
		},
	}
	return domain.SourdoughRecipeEntity{
		RecipeEntity: domain.RecipeEntity{
			Id:          id,
			Name:        fmt.Sprintf("name-%s", id.String()),
			Description: fmt.Sprintf("description-%s", id.String()),
			Flour: []domain.FlourAmount{
				{
					FlourEntity: flour,
					Amount:      1000,
				},
			},
			Water: []domain.BakerAmount{
				{
					Amount:          750,
					BakerPercentage: 0.75,
					Name:            "Water 1",
				},
				{
					Amount:          50,
					BakerPercentage: 0.05,
					Name:            "Water 2",
				},
				{
					Amount:          50,
					BakerPercentage: 0.05,
					Name:            "Water for levain",
				},
			},
			AdditionalIngredients: []domain.BakerAmount{
				{
					Amount:          20,
					BakerPercentage: 0.02,
					Name:            "Salt",
				},
			},
			Details: domain.RecipeDetails{
				Flour: domain.BakerAmount{
					Amount:          1000,
					BakerPercentage: 1,
					Name:            "Total FlourDto",
				},
				Water: domain.BakerAmount{
					Amount:          800,
					BakerPercentage: 0.8,
					Name:            "Total Water",
				},
				Levain: domain.BakerAmount{
					Amount:          20,
					BakerPercentage: 0.02,
					Name:            "Total Levain",
				},
				AdditionalIngredients: domain.BakerAmount{
					Amount:          20,
					BakerPercentage: 0.02,
					Name:            "Total Additional Ingredients",
				},
				TotalWeight: 1840,
			},
			NutritionFacts: map[string]domain.NutritionFacts{
				"100g": {
					Calories: 100,
					Fat:      1,
					Carbs:    1,
					Protein:  2.5,
					Fiber:    1,
				},
				"1 loaf": {
					Calories: 200,
					Fat:      2,
					Carbs:    2,
					Protein:  5,
					Fiber:    2,
				},
			},
			CreatedAt: time.Now().Truncate(time.Second).UTC(),
			Yield: domain.RecipeYield{
				Unit:   "loaf",
				Amount: 2,
			},
		},
		Levain: domain.SourdoughLevainAgent{
			Starter: domain.BakerAmount{
				Amount:          66.67,
				BakerPercentage: 0.0667,
				Name:            "Starter",
			},
			Flour: []domain.FlourAmount{
				{
					FlourEntity: flour,
					Amount:      66.67,
				},
			},
			Water: domain.BakerAmount{
				Amount:          66.67,
				BakerPercentage: 0.0667,
				Name:            "Water for levain",
			},
		},
	}
}
