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
	var res ErrorResponse
	if errors.As(err, &res) {
		return res.StatusCode
	}
	return http.StatusInternalServerError
}

// ErrDataValidation is used when validation by `ozzo-validation` returns an error.
// Error value from `ozzo-validation` can be marshaled into key-value JSON object.
func ErrDataValidation(errs validation.Errors) error {
	resData := make(map[string]string, len(errs))
	for k, v := range errs {
		if err, ok := v.(validation.Error); ok {
			// Capitalize first letter of the error messages and add period at the end.
			errMsg := err.Message()
			resData[k] = string(errMsg[0]-32) + errMsg[1:] + "."
		} else {
			// Some errors are not of type `validation.Error`
			// which can happen due to code bug instead of user input.
			// E.g. "cannot get the length of struct" when using length rules on a struct.
			return fmt.Errorf("validate field '%s': %w", k, v) // INTERNAL SERVER ERROR
		}
	}
	return ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "Input verification failed.",
		Data:       resData,
	}
}

func ErrPermission() error {
	// TODO: Should return 404 instead of 403?
	return Error(http.StatusForbidden, "You are not authorized to perform this request.")
}

func ErrIDNotFound(id uuid.UUID) error {
	return Errorf(http.StatusNotFound, "ID '%s' not found.", id)
}

func ErrConflict(key string, value string) error {
	return Errorf(http.StatusConflict, "%s '%s' already exists.", key, value)
}
