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
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func Error(statusCode int, message any) error {
	return &APIError{statusCode, message}
}

func Errorf(statusCode int, format string, args ...interface{}) error {
	return &APIError{statusCode, fmt.Sprintf(format, args...)}
}

func ErrorStatusCode(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return http.StatusInternalServerError
}

// ErrInvalidJSON is used when JSON decoder failed to parse the request body.
func ErrInvalidJSON() error {
	return Error(http.StatusBadRequest, "invalid request body")
}

// ErrInvalidUUID is used when the request param is not a valid UUID.
func ErrInvalidUUID(invalidUUID string) error {
	return Errorf(http.StatusBadRequest, "invalid UUID '%s'", invalidUUID)
}

// ErrValidation is used when validation by `ozzo-validation` returns an error.
// Error message from `ozzo-validation` can be marshaled into key-value JSON object.
func ErrValidation(err error) error {
	return Error(http.StatusBadRequest, err)
}

func ErrNotFound(id uuid.UUID) error {
	return Errorf(http.StatusBadRequest, "ID '%s' not found", id)
}

func ErrConflict(key string, value string) error {
	return Errorf(http.StatusBadRequest, "%s '%s' already exists", key, value)
}
