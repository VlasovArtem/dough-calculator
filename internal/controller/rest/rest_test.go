//go:build integration

package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalErrors "dough-calculator/internal/errors"

	"github.com/stretchr/testify/assert"
)

func TestWriteServiceError(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	testErr := internalErrors.NewServiceError(500, -1, "Service error message", "Service error details")

	handlerServiceError(rec, req, testErr.(*internalErrors.ServiceError))

	assert.Equal(t, 500, rec.Code)
	expected := `{"error_code":-1,"error_message":"Service error message","error_details":"Service error details"}`
	assert.JSONEq(t, expected, rec.Body.String())
}

func TestWriteError(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	testErr := internalErrors.NewServiceError(500, -1, "Service error message", "Service error details")

	HandlerError(rec, req, testErr)

	assert.Equal(t, 500, rec.Code)
	expected := `{"error_code":-1,"error_message":"Service error message","error_details":"Service error details"}`
	assert.JSONEq(t, expected, rec.Body.String())
}

func TestWriteError_WithErrorNotService(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	testErr := errors.New("test error")

	HandlerError(rec, req, testErr)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error_code":-1, "error_details":"test error", "error_message":"internal server error"}`, rec.Body.String())
}

func TestWriteError_WithNoError(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	HandlerError(rec, req, nil)

	assert.Equal(t, 200, rec.Code)
}
