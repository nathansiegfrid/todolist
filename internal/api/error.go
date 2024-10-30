package api

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	StatusCode int
	Message    string
	Data       any
}

// Error implements the `error` interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

func Error(statusCode int, message string) error {
	return ErrorResponse{StatusCode: statusCode, Message: message}
}

func Errorf(statusCode int, format string, args ...any) error {
	return ErrorResponse{StatusCode: statusCode, Message: fmt.Sprintf(format, args...)}
}

func ErrorStatusCode(err error) int {
	var apiErr ErrorResponse
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return http.StatusInternalServerError
}

// ErrDataValidation is used when validation by `ozzo-validation` returns an error.
// Error value from `ozzo-validation` can be marshaled into key-value JSON object.
func ErrDataValidation(err error) error {
	var errs validation.Errors
	if errors.As(err, &errs) {
		data := make(map[string]string, len(errs))
		// Capitalize first letter of the error messages and add period at the end.
		for k, v := range errs {
			errMsg := v.Error()
			if len(errMsg) > 0 {
				data[k] = string(errMsg[0]-32) + errMsg[1:] + "."
			}
		}
		return ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Input verification failed.",
			Data:       data,
		}
	}
	return err // Internal server error.
}

func ErrPermission() error {
	return Error(http.StatusForbidden, "Permission denied.")
}

func ErrIDNotFound(id uuid.UUID) error {
	return Errorf(http.StatusNotFound, "ID '%s' not found.", id)
}

func ErrConflict(key string, value string) error {
	return Errorf(http.StatusConflict, "%s '%s' already exists.", key, value)
}
