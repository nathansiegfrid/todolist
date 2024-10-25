package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
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

// MethodHandler maps HTTP methods to handler functions.
type MethodHandler map[string]http.HandlerFunc

func (hmap MethodHandler) HandlerFunc() http.HandlerFunc {
	// Define HTTP methods that can be used as keys in MethodHandler.
	// Unknown keys will return false by default.
	validKeys := map[string]bool{
		"GET":     true,
		"HEAD":    true,
		"POST":    true,
		"PUT":     true,
		"PATCH":   true,
		"DELETE":  true,
		"CONNECT": false,
		"OPTIONS": false,
		"TRACE":   false,
	}

	// Get list of allowed methods.
	methods := make([]string, 0, len(hmap)+1)
	methods = append(methods, "OPTIONS")
	for k, v := range hmap {
		// Check if method is valid and handler is not nil.
		if validKeys[k] && v != nil {
			methods = append(methods, k)
		} else {
			delete(hmap, k)
		}
	}
	allowHeaderValue := strings.Join(methods, ", ")

	// Implement OPTIONS method to handle preflight requests.
	hmap["OPTIONS"] = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", allowHeaderValue)
		w.WriteHeader(http.StatusNoContent)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		serveHTTP := hmap[r.Method]
		if serveHTTP == nil {
			// If no handler found, respond with 405 Method Not Allowed.
			w.Header().Set("Allow", allowHeaderValue)
			MethodNotAllowed(w, r)
			return
		}
		serveHTTP(w, r)
	}
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	err := Error(http.StatusNotFound, "Resource not found.")
	WriteError(w, err)
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	err := Error(http.StatusMethodNotAllowed, "Method not allowed.")
	WriteError(w, err)
}
