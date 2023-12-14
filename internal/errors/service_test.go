package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceError(t *testing.T) {
	err := NewServiceError(400, 123, "Test error", "Test details")

	assert.Equal(t, &ServiceError{
		ResponseCode: 400,
		Code:         123,
		Message:      "Test error",
		Details:      "Test details",
	}, err)
}

func TestNewServiceErrorWrap(t *testing.T) {
	cause := fmt.Errorf("causal error")
	err := NewServiceErrorWrap(400, 123, "Test error", cause)

	assert.Equal(t, &ServiceError{
		ResponseCode: 400,
		Code:         123,
		Message:      "Test error",
		Details:      "causal error",
		causedBy:     cause,
	}, err)
}

func TestError(t *testing.T) {
	se := &ServiceError{
		ResponseCode: 400,
		Code:         123,
		Message:      "Test error",
		Details:      "Test details",
	}

	assert.Equal(t, "error code: 123, error message: Test error, error details: Test details", se.Error())
}

func TestCause(t *testing.T) {
	cause := fmt.Errorf("causal error")
	se := &ServiceError{
		ResponseCode: 400,
		Code:         123,
		Message:      "Test error",
		Details:      "Test details",
		causedBy:     cause,
	}

	assert.Equal(t, cause, se.Cause())
}

func TestIsServiceError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ServiceError",
			err:      &ServiceError{},
			expected: true,
		},
		{
			name:     "NotServiceError",
			err:      fmt.Errorf("causal error"),
			expected: false,
		},
		{
			name:     "Nil",
			err:      nil,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsServiceError(test.err)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name          string
		provider      func() error
		expectedError *ServiceError
	}{
		{
			name: "BadRequestError",
			provider: func() error {
				return NewBadRequestError(123, "Test error", "Test details")
			},
			expectedError: &ServiceError{
				ResponseCode: http.StatusBadRequest,
				Code:         123,
				Message:      "Test error",
				Details:      "Test details",
			},
		},
		{
			name: "BadRequestErrorf",
			provider: func() error {
				return NewBadRequestErrorf(123, "Test error", "Test details %s", "test")
			},
			expectedError: &ServiceError{
				ResponseCode: http.StatusBadRequest,
				Code:         123,
				Message:      "Test error",
				Details:      "Test details test",
			},
		},
		{
			name: "InternalServerError",
			provider: func() error {
				return NewInternalServerError("Test error", "Test details")
			},
			expectedError: &ServiceError{
				ResponseCode: http.StatusInternalServerError,
				Code:         -1,
				Message:      "Test error",
				Details:      "Test details",
			},
		},
		{
			name: "InternalServerErrorWrap",
			provider: func() error {
				return NewInternalServerErrorWrap(fmt.Errorf("causal error"), "Test error")
			},
			expectedError: &ServiceError{
				ResponseCode: http.StatusInternalServerError,
				Code:         -1,
				Message:      "Test error",
				Details:      "causal error",
				causedBy:     fmt.Errorf("causal error"),
			},
		},
		{
			name: "UnauthorizedServerError",
			provider: func() error {
				return NewUnauthorizedServerError("Test error", "Test details")
			},
			expectedError: &ServiceError{
				ResponseCode: http.StatusUnauthorized,
				Code:         401,
				Message:      "Test error",
				Details:      "Test details",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedError, tt.provider())
		})
	}
}
