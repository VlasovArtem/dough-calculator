package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"dough-calculator/internal/domain"
	"dough-calculator/internal/domain/mocks"
	internalErrors "dough-calculator/internal/errors"
	"dough-calculator/internal/test"
)

func TestSourdoughRecipeScaleHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SourdoughRecipeScaleHandlerTestSuite))
}

type SourdoughRecipeScaleHandlerTestSuite struct {
	test.GoMockTestSuite

	service *mocks.MockSourdoughRecipeScaleService

	target domain.SourdoughRecipeScaleHandler
}

func (suite *SourdoughRecipeScaleHandlerTestSuite) SetupTest() {
	suite.GoMockTestSuite.SetupTest()

	suite.service = mocks.NewMockSourdoughRecipeScaleService(suite.MockCtrl)

	suite.target = test.Must(func() (domain.SourdoughRecipeScaleHandler, error) {
		return NewSourdoughRecipeScaleHandler(suite.service)
	})
}

func (suite *SourdoughRecipeScaleHandlerTestSuite) TestScale() {
	id := uuid.New()
	request := domain.SourdoughRecipeScaleRequestDto{
		FinalDoughWeight: 500,
	}

	suite.service.EXPECT().
		Scale(gomock.Any(), id, request).
		Return(createSourdoughRecipe(), nil)

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/scale/{id}", suite.target.Scale())

	req, err := http.NewRequest("POST", fmt.Sprintf("/scale/%s", id), buffer)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	test.VerifyRestResponseWithTestFile(suite.T(), resp, http.StatusOK, "testdata/sourdough_recipe_response.json")
}

func (suite *SourdoughRecipeScaleHandlerTestSuite) TestScale_WithErrorOnScale() {
	id := uuid.New()
	request := domain.SourdoughRecipeScaleRequestDto{
		FinalDoughWeight: 500,
	}

	suite.service.EXPECT().
		Scale(gomock.Any(), id, request).
		Return(domain.SourdoughRecipeDto{}, internalErrors.NewBadRequestErrorf(123, "error", "error %s", "'test'"))

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/scale/{id}", suite.target.Scale())

	req, err := http.NewRequest("POST", fmt.Sprintf("/scale/%s", id), buffer)
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

func (suite *SourdoughRecipeScaleHandlerTestSuite) TestScale_WithInvalidBody() {
	id := uuid.New()

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode("invalid")
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/scale/{id}", suite.target.Scale())

	req, err := http.NewRequest("POST", fmt.Sprintf("/scale/%s", id), buffer)
	suite.Require().NoError(err)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	expectedBodyJson :=
		`{
			"error_code": -1,
			"error_details": "error while decoding request body: json: cannot unmarshal string into Go value of type domain.SourdoughRecipeScaleRequestDto",
			"error_message": "internal server error"
		}`
	test.VerifyRestResponse(suite.T(), resp, http.StatusInternalServerError, expectedBodyJson)
}

func (suite *SourdoughRecipeScaleHandlerTestSuite) TestScale_WithInvalidId() {
	request := domain.SourdoughRecipeScaleRequestDto{
		FinalDoughWeight: 500,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/scale/{id}", suite.target.Scale())

	req, err := http.NewRequest("POST", "/scale/invalid", buffer)
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

func (suite *SourdoughRecipeScaleHandlerTestSuite) TestScale_WithoutRequiredIdParam() {
	request := domain.SourdoughRecipeScaleRequestDto{
		FinalDoughWeight: 500,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buffer).Encode(request)
	suite.Require().NoError(err)

	router := chi.NewRouter()
	router.
		Post("/scale", suite.target.Scale())

	req, err := http.NewRequest("POST", "/scale", buffer)
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

func TestNewSourdoughRecipeScaleHandler_WithNilService(t *testing.T) {
	_, err := NewSourdoughRecipeScaleHandler(nil)

	assert.ErrorContains(t, err, "service is nil")
}
