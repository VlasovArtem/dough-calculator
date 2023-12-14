package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ServiceError struct {
	ResponseCode int    `json:"-"`
	Code         int    `json:"error_code"`
	Message      string `json:"error_message"`
	Details      string `json:"error_details"`
	causedBy     error
}

func (err *ServiceError) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s, error details: %s", err.Code, err.Message, err.Details)
}

func (err *ServiceError) Cause() error {
	return err.causedBy
}

func NewServiceError(responseCode, serviceErrorCode int, message, details string) error {
	return &ServiceError{
		ResponseCode: responseCode,
		Code:         serviceErrorCode,
		Message:      message,
		Details:      details,
	}
}

func NewServiceErrorWrap(responseCode, serviceErrorCode int, message string, err error) error {
	return &ServiceError{
		ResponseCode: responseCode,
		Code:         serviceErrorCode,
		Message:      message,
		Details:      err.Error(),
		causedBy:     err,
	}
}

func NewBadRequestError(serviceErrorCode int, message, details string) error {
	return NewServiceError(http.StatusBadRequest, serviceErrorCode, message, details)
}

func NewBadRequestErrorf(serviceErrorCode int, message string, details string, args ...any) error {
	return NewServiceError(http.StatusBadRequest, serviceErrorCode, message, fmt.Sprintf(details, args...))
}

func NewInternalServerError(message, details string) error {
	return NewServiceError(http.StatusInternalServerError, -1, message, details)
}

func NewInternalServerErrorWrap(err error, message string) error {
	return NewServiceErrorWrap(http.StatusInternalServerError, -1, message, err)
}

func NewUnauthorizedServerError(message, details string) error {
	return NewServiceError(http.StatusUnauthorized, 401, message, details)
}

func IsServiceError(err error) bool {
	if err == nil {
		return false
	}

	var serviceError *ServiceError
	ok := errors.As(err, &serviceError)
	return ok
}

func DbIsNilError() error {
	return errors.New("db is nil")
}
