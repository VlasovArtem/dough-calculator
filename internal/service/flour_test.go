package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	internalErrors "dough-calculator/internal/errors"
	"dough-calculator/internal/test"
)

func TestFlourServiceTestSuite(t *testing.T) {
	suite.Run(t, new(FlourServiceTestSuite))
}

type FlourServiceTestSuite struct {
	test.GoMockTestSuite

	ctx        context.Context
	repository *mocks.MockFlourRepository

	target domain.FlourService
}

func (suite *FlourServiceTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.ctx = context.Background()
	suite.repository = mocks.NewMockFlourRepository(suite.MockCtrl)

	suite.target = test.Must(func() (domain.FlourService, error) {
		return NewFlourService(suite.repository)
	})
}

func (suite *FlourServiceTestSuite) TestCreate() {
	createRequest := suite.createRequest()

	var savedEntity domain.FlourEntity
	suite.repository.EXPECT().Create(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, entity domain.FlourEntity) (domain.FlourEntity, error) {
			savedEntity = entity
			return entity, nil
		})

	actualDto, err := suite.target.Create(suite.ctx, createRequest)

	suite.NoError(err)
	suite.Equal(savedEntity.ToDto(), actualDto)
}

func (suite *FlourServiceTestSuite) TestCreate_WithError() {
	createRequest := suite.createRequest()

	suite.repository.EXPECT().Create(suite.ctx, gomock.Any()).
		Return(domain.FlourEntity{}, assert.AnError)

	_, err := suite.target.Create(suite.ctx, createRequest)

	suite.ErrorContains(err, "failed to create flour")
}

func (suite *FlourServiceTestSuite) createRequest() domain.CreateFlourRequest {
	return domain.CreateFlourRequest{
		FlourType:   "Test FlourDto",
		Name:        "Test Name",
		Description: "Test Description",
		NutritionFacts: domain.NutritionFactsDto{
			Calories: 100,
			Fat:      1,
			Carbs:    1,
			Protein:  1,
			Fiber:    1,
		},
	}
}

func (suite *FlourServiceTestSuite) TestFindById() {
	entity := suite.createEntity()

	suite.repository.EXPECT().FindById(suite.ctx, entity.Id).
		Return(entity, nil)

	actualDto, err := suite.target.FindById(suite.ctx, entity.Id)

	suite.NoError(err)
	suite.Equal(entity.ToDto(), actualDto)
}

func (suite *FlourServiceTestSuite) TestFindById_WithError() {
	entity := suite.createEntity()

	tests := []struct {
		name                string
		errorFromRepository error
		expectedError       error
	}{
		{
			name:                "with basic error",
			errorFromRepository: assert.AnError,
			expectedError:       internalErrors.NewInternalServerErrorWrap(assert.AnError, "failed to find flour by id"),
		},
		{
			name:                "with document not found error",
			errorFromRepository: mongo.ErrNoDocuments,
			expectedError:       internalErrors.FlourByIdNotFound(entity.Id),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.repository.EXPECT().
				FindById(suite.ctx, entity.Id).
				Return(entity, tt.errorFromRepository)

			_, err := suite.target.FindById(suite.ctx, entity.Id)

			suite.Equal(tt.expectedError, err)
		})
	}
}

func (suite *FlourServiceTestSuite) TestFind() {
	entity := suite.createEntity()

	suite.repository.EXPECT().Find(suite.ctx, 0, 10).
		Return([]domain.FlourEntity{entity}, nil)

	actualDto, err := suite.target.Find(suite.ctx, 0, 10)

	suite.NoError(err)
	suite.Equal([]domain.FlourDto{entity.ToDto()}, actualDto)
}

func (suite *FlourServiceTestSuite) TestFind_WithError() {
	suite.repository.EXPECT().Find(suite.ctx, 0, 10).
		Return([]domain.FlourEntity{}, assert.AnError)

	_, err := suite.target.Find(suite.ctx, 0, 10)

	suite.ErrorContains(err, "failed to find flours")
}

func (suite *FlourServiceTestSuite) TestSearchByName() {
	entity := suite.createEntity()

	suite.repository.EXPECT().SearchByName(suite.ctx, entity.Name).
		Return([]domain.FlourEntity{entity}, nil)

	actualDto, err := suite.target.SearchByName(suite.ctx, entity.Name)

	suite.NoError(err)
	suite.Equal([]domain.FlourDto{entity.ToDto()}, actualDto)
}

func (suite *FlourServiceTestSuite) TestSearchByName_WithError() {
	entity := suite.createEntity()

	suite.repository.EXPECT().SearchByName(suite.ctx, entity.Name).
		Return([]domain.FlourEntity{}, assert.AnError)

	_, err := suite.target.SearchByName(suite.ctx, entity.Name)

	suite.ErrorContains(err, "failed to search flours by name")
}

func (suite *FlourServiceTestSuite) createEntity() domain.FlourEntity {
	return domain.FlourEntity{
		Id:          test.FirstId,
		FlourType:   "Test FlourDto",
		Name:        "Test Name",
		Description: "Test Description",
		NutritionFacts: domain.NutritionFacts{
			Calories: 100,
			Fat:      1,
			Carbs:    1,
			Protein:  1,
			Fiber:    1,
		},
	}
}

func TestNewFlourService_WithNilRepository(t *testing.T) {
	_, err := NewFlourService(nil)

	assert.EqualError(t, err, "repository is nil")
}
