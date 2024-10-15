package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/schema"
)

// ReadURLQuery maps URL query into struct using `schema` tags.
// Supports primitive types, time.Time, and uuid.UUID.
func ReadURLQuery[T any](r *http.Request) (*T, error) {
	dec := schema.NewDecoder()
	// Add custom converter because default decoder doesn't support time.Time.
	dec.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		t, err := time.Parse(time.DateOnly, value) // Format: YYYY-MM-DD.
		if err != nil {
			return reflect.Value{} // Zero value represents invalid Value and "Invalid" kind.
		}
		return reflect.ValueOf(t)
	})

	dst := new(T)
	err := dec.Decode(dst, r.URL.Query())
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func ReadJSON[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

type responseBody struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func write(w http.ResponseWriter, statusCode int, body *responseBody) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

func WriteOK(w http.ResponseWriter) error {
	return write(w, http.StatusOK, &responseBody{Status: "SUCCESS"})
}

func WriteJSON(w http.ResponseWriter, data any) error {
	return write(w, http.StatusOK, &responseBody{Status: "SUCCESS", Data: data})
}

func WriteError(w http.ResponseWriter, err error) error {
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode != http.StatusInternalServerError {
		return write(w, apiErr.StatusCode, &responseBody{
			Status:  "FAIL",
			Message: apiErr.Message,
			Data:    apiErr.Data,
		})
	}

	return write(w, http.StatusInternalServerError, &responseBody{
		Status:  "ERROR",
		Message: "Unexpected error. We've noted the issue. Please try again later.", // Log the error.
	})
}
