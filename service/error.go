package service

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func Errorf(statusCode int, format string, args ...any) error {
	return &APIError{statusCode, fmt.Sprintf(format, args...)}
}

func ErrorStatusCode(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return http.StatusInternalServerError
}

// ErrInvalidID is used when the request param is not a valid ID.
func ErrInvalidID(id string) error {
	return Errorf(http.StatusBadRequest, "invalid ID '%s'", id)
}

// ErrInvalidJSON is used when JSON decoder failed to parse the request body.
func ErrInvalidJSON() error {
	return Error(http.StatusBadRequest, "invalid request body")
}

// ErrValidation is used when validation by `ozzo-validation` returns an error.
// Error message from `ozzo-validation` can be marshaled into key-value JSON object.
func ErrValidation(errs validation.Errors) error {
	return Error(http.StatusBadRequest, errs)
}

// ErrPermission is used when the user does not have enough authorization to do the request.
func ErrPermission() error {
	return Error(http.StatusForbidden, "permission denied")
}

func ErrNotFound(id uuid.UUID) error {
	return Errorf(http.StatusNotFound, "ID '%s' not found", id)
}

func ErrConflict(key string, value string) error {
	return Errorf(http.StatusConflict, "%s '%s' already exists", key, value)
}
