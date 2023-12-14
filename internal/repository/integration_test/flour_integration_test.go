//go:build integration && docker

package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/repository"
	"dough-calculator/internal/test"
)

func TestFlourRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &FlourRepositoryTestSuite{
		MongoDBServiceDockerIntegrationTestSuite: test.NewMongoDBServiceDockerIntegrationTestSuite(dockerStarter),
	})
}

type FlourRepositoryTestSuite struct {
	test.MongoDBServiceDockerIntegrationTestSuite

	mockMongoDbService *mocks.MockMongoDBService

	target domain.FlourRepository
}

func (suite *FlourRepositoryTestSuite) SetupSuite() {
	suite.MongoDBServiceDockerIntegrationTestSuite.SetupSuite()

	suite.target = test.Must(func() (domain.FlourRepository, error) {
		return repository.NewFlourRepository(suite.Stub)
	})
}

func (suite *FlourRepositoryTestSuite) AfterTest(suiteName, testName string) {
	err := suite.Drop(repository.FlourDatabase, repository.FlourCollection)
	suite.Require().NoError(err)
}

func (suite *FlourRepositoryTestSuite) TestCreate() {
	expected := generateFlourEntity()

	actual, err := suite.target.Create(context.Background(), expected)

	suite.NoError(err)
	suite.Equal(expected, actual)

	var saved domain.FlourEntity

	err = suite.MStub().MustGetCollection(repository.FlourDatabase, repository.FlourCollection).
		FindOne(context.Background(), bson.D{{"_id", expected.Id}}).
		Decode(&saved)

	suite.Equal(expected, saved)
}

func (suite *FlourRepositoryTestSuite) TestCreate_WithEntityExists_ShouldReturnError() {
	expected := generateFlourEntity()

	actual, err := suite.target.Create(context.Background(), expected)

	suite.NoError(err)
	suite.Equal(expected, actual)

	_, err = suite.target.Create(context.Background(), expected)

	suite.ErrorContains(err, "failed to insert flour")
}

func (suite *FlourRepositoryTestSuite) TestGetById() {
	expected := generateFlourEntity()

	_, err := suite.target.Create(context.Background(), expected)
	suite.NoError(err)

	actual, err := suite.target.FindById(context.Background(), expected.Id)

	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *FlourRepositoryTestSuite) TestGetById_WithEntityNotFound_ShouldReturnError() {
	expected := generateFlourEntity()

	_, err := suite.target.FindById(context.Background(), expected.Id)

	suite.ErrorContains(err, "failed to get flour by id")
}

func (suite *FlourRepositoryTestSuite) TestFind() {
	first := generateFlourEntity()
	_, err := suite.target.Create(context.Background(), first)
	suite.Require().NoError(err)
	second := generateFlourEntity()
	_, err = suite.target.Create(context.Background(), second)
	suite.Require().NoError(err)

	actual, err := suite.target.Find(context.Background(), 0, 1)

	suite.NoError(err)
	suite.Contains(actual, first)
}

func (suite *FlourRepositoryTestSuite) TestFind_WithEmptyData_ShouldReturnNil() {
	actual, err := suite.target.Find(context.Background(), 1, 0)

	suite.NoError(err)
	suite.Nil(actual)
}

func (suite *FlourRepositoryTestSuite) TestFindByName() {
	entity := generateFlourEntity()
	_, err := suite.target.Create(context.Background(), entity)
	suite.Require().NoError(err)

	actual, err := suite.target.SearchByName(context.Background(), entity.Name)

	suite.NoError(err)
	suite.Equal([]domain.FlourEntity{entity}, actual)
}

func (suite *FlourRepositoryTestSuite) TestFindByName_WithEntityNotExists_ShouldReturnEmptyEntity() {
	actual, err := suite.target.SearchByName(context.Background(), "missing")

	suite.NoError(err)
	suite.Nil(actual)
}

func generateFlourEntity() domain.FlourEntity {
	id := uuid.New()

	return domain.FlourEntity{
		Id:          id,
		FlourType:   fmt.Sprintf("test %s", id.String()),
		Name:        fmt.Sprintf("test flour %s", id.String()),
		Description: fmt.Sprintf("test flour description %s", id.String()),
		NutritionFacts: domain.NutritionFacts{
			Calories: 100,
			Fat:      1,
			Carbs:    1,
			Protein:  2.5,
			Fiber:    1,
		},
	}
}
