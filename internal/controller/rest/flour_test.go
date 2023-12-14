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

func TestFlourHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FlourHandlerTestSuite))
}

type FlourHandlerTestSuite struct {
	test.GoMockTestSuite

	service *mocks.MockFlourService

	target domain.FlourHandler
}

func (suite *FlourHandlerTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.service = mocks.NewMockFlourService(suite.MockCtrl)

	suite.target = test.Must(func() (domain.FlourHandler, error) {
		return NewFlourHandler(suite.service)
	})
}

func (suite *FlourHandlerTestSuite) TestCreateFlour() {
	request := generateCreateFlourRequest()

	suite.service.EXPECT().
		Create(gomock.Any(), request).
		Return(createFlour(), nil)

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

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusCreated, "testdata/flour_response.json")
}

func (suite *FlourHandlerTestSuite) TestCreateFlour_WithInvalidRequest() {
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

func (suite *FlourHandlerTestSuite) TestCreateFlour_WithErrorOnCreate() {
	request := generateCreateFlourRequest()

	suite.service.EXPECT().
		Create(gomock.Any(), request).
		Return(domain.FlourDto{}, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

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

func (suite *FlourHandlerTestSuite) TestFindFlourById() {
	flour := createFlour()

	suite.service.EXPECT().FindById(gomock.Any(), flour.Id).
		Return(flour, nil)

	router := chi.NewRouter()
	router.
		Get("/flour/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", fmt.Sprintf("/flour/%s", flour.Id), nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/flour_response.json")
}

func (suite *FlourHandlerTestSuite) TestFindFlourById_WithoutParam() {
	router := chi.NewRouter()
	router.
		Get("/flour", suite.target.FindById())

	req, err := http.NewRequest("GET", "/flour", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
		"error_code": 20001,
		"error_details": "id is required",
		"error_message": "id is required"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *FlourHandlerTestSuite) TestFindFlourById_WithInvalidIdParam() {
	router := chi.NewRouter()
	router.
		Get("/flour/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", "/flour/invalid", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
		"error_code": 20002,
		"error_details": "id is not valid",
		"error_message": "id is not valid"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusBadRequest, expectedBodyJson)
}

func (suite *FlourHandlerTestSuite) TestFindFlourById_WithErrorOnFindById() {
	flour := createFlour()

	suite.service.EXPECT().FindById(gomock.Any(), flour.Id).
		Return(domain.FlourDto{}, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	router := chi.NewRouter()
	router.
		Get("/flour/{id}", suite.target.FindById())

	req, err := http.NewRequest("GET", fmt.Sprintf("/flour/%s", flour.Id), nil)
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

func (suite *FlourHandlerTestSuite) FindFlour() {
	flours := []domain.FlourDto{createFlour()}

	suite.service.EXPECT().Find(gomock.Any(), 1, 10).Return(flours, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(PageInput{})).
		Get("/find", suite.target.Find())

	req, err := http.NewRequest("GET", "/find?offset=1&limit=10", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/flours_response.json")
}

func (suite *FlourHandlerTestSuite) TestFindFlour_WithDefaultParameters() {
	flours := []domain.FlourDto{createFlour()}

	suite.service.EXPECT().Find(gomock.Any(), 0, 25).Return(flours, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(PageInput{})).
		Get("/", suite.target.Find())

	req, err := http.NewRequest("GET", "/", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/flours_response.json")
}

func (suite *FlourHandlerTestSuite) TestFindFlour_WithErrorOnFind() {
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

func (suite *FlourHandlerTestSuite) TestSearchFlourByName() {
	flours := []domain.FlourDto{createFlour()}

	suite.service.EXPECT().SearchByName(gomock.Any(), "test name").
		Return(flours, nil)

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(SearchFlourInput{})).
		Get("/search", suite.target.Search())

	req, err := http.NewRequest("GET", "/search?name=test%20name", nil)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/flours_response.json")
}

func (suite *FlourHandlerTestSuite) TestSearchFlourByName_WithErrorOnSearch() {
	suite.service.EXPECT().SearchByName(gomock.Any(), "test name").
		Return(nil, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	router := chi.NewRouter()
	router.
		With(httpin.NewInput(SearchFlourInput{})).
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

func TestNewFlourHandler_WithNilService(t *testing.T) {
	_, err := NewFlourHandler(nil)

	assert.ErrorContains(t, err, "service is nil")
}

func generateCreateFlourRequest() domain.CreateFlourRequest {
	return domain.CreateFlourRequest{
		FlourType:   "Whole Wheat",
		Name:        "Whole Wheat Flour",
		Description: "Whole grain flour milled from red wheat berries",
		NutritionFacts: domain.NutritionFactsDto{
			Calories: 100,
			Fat:      1,
			Carbs:    21,
			Protein:  4,
			Fiber:    3,
		},
	}
}

func createFlour() domain.FlourDto {
	return domain.FlourDto{
		Id:          test.FirstId,
		FlourType:   "Whole Wheat",
		Name:        "Whole Wheat Flour",
		Description: "Whole grain flour milled from red wheat berries",
		NutritionFacts: domain.NutritionFactsDto{
			Calories: 100,
			Fat:      1,
			Carbs:    21,
			Protein:  4,
			Fiber:    3,
		},
	}
}
