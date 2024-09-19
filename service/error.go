package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type APIError struct {
	StatusCode int
	Message    any
}

// Error implements the `error` interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("error %d: %s", e.StatusCode, e.Message)
}

func Error(statusCode int, message any) *APIError {
	return &APIError{statusCode, message}
}

func Errorf(statusCode int, format string, args ...interface{}) *APIError {
	return &APIError{statusCode, fmt.Sprintf(format, args...)}
}

func StatusCode(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return http.StatusInternalServerError
}

// ErrInvalidJSON is used when JSON decoder failed to parse the request body.
func ErrInvalidJSON() *APIError {
	return Error(http.StatusBadRequest, "invalid request body: JSON required")
}

// ErrInvalidUUID is used when the request param is not a valid UUID.
func ErrInvalidUUID(invalidUUID string) *APIError {
	return Errorf(http.StatusBadRequest, "invalid UUID: '%s'", invalidUUID)
}

// ErrValidation is used when validation by `ozzo-validation` returns an error.
// Error message from `ozzo-validation` can be marshaled into key-value JSON object.
func ErrValidation(err error) *APIError {
	return Error(http.StatusBadRequest, err)
}

func ErrNotFound(id uuid.UUID) *APIError {
	return Errorf(http.StatusBadRequest, "ID not found: '%s'", id)
}

func ErrConflict(key string, value string) *APIError {
	return Errorf(http.StatusBadRequest, "%s already exists: '%s'", key, value)
}
