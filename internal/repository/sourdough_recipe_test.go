package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeRepositoryTestSuite))
}

type SourdoughRecipeRepositoryTestSuite struct {
	test.GoMockTestSuite

	mongoDBService *mocks.MockMongoDBService

	target *sourdoughRecipeRepository
}

func (suite *SourdoughRecipeRepositoryTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.mongoDBService = mocks.NewMockMongoDBService(suite.MockCtrl)

	suite.target = &sourdoughRecipeRepository{
		mongoDBService: suite.mongoDBService,
	}
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestNewSourdoughRecipeRepository_WithError() {
	tests := []struct {
		name           string
		mongoDBService domain.MongoDBService
		errorMsg       string
	}{
		{
			name:           "mongoDBService is nil",
			mongoDBService: nil,
			errorMsg:       "service cannot be nil",
		},
		{
			name: "mongoDBService.GetCollection returns error",
			mongoDBService: func() domain.MongoDBService {
				suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
					Return(nil, assert.AnError)

				return suite.mongoDBService
			}(),
			errorMsg: "failed to get collection",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			repository, err := NewSourdoughRecipeRepository(tt.mongoDBService)

			suite.ErrorContains(err, tt.errorMsg)
			suite.Nil(repository)
		})
	}
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestGetCollection_WithError() {
	suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
		Return(nil, assert.AnError)

	collection, err := suite.target.getCollection()

	suite.ErrorContains(err, "failed to get collection")
	suite.Nil(collection)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestCreate_WithErrorOnGetCollection() {
	suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
		Return(nil, assert.AnError)

	recipe := domain.SourdoughRecipeEntity{}

	entity, err := suite.target.Create(context.Background(), recipe)

	suite.ErrorContains(err, "failed to get collection")
	suite.Equal(domain.SourdoughRecipeEntity{}, entity)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestGetById_WithErrorOnGetCollection() {
	suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
		Return(nil, assert.AnError)

	entity, err := suite.target.GetById(context.Background(), uuid.UUID{})

	suite.ErrorContains(err, "failed to get collection")
	suite.Equal(domain.SourdoughRecipeEntity{}, entity)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestFind_WithErrorOnGetCollection() {
	suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
		Return(nil, assert.AnError)

	entities, err := suite.target.Find(context.Background(), 0, 1)

	suite.ErrorContains(err, "failed to get collection")
	suite.Nil(entities)
}

func (suite *SourdoughRecipeRepositoryTestSuite) TestSearchByName_WithErrorOnGetCollection() {
	suite.mongoDBService.EXPECT().GetCollection(SourdoughRecipeDatabase, SourdoughRecipeCollection).
		Return(nil, assert.AnError)

	entities, err := suite.target.SearchByName(context.Background(), "")

	suite.ErrorContains(err, "failed to get collection")
	suite.Nil(entities)
}
