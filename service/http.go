package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
)

// ReadURLQuery maps URL query into struct using `schema` tags.
func ReadURLQuery[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := schema.NewDecoder().Decode(dst, r.URL.Query())
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

type response struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func write(w http.ResponseWriter, statusCode int, response response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}

func WriteOK(w http.ResponseWriter) error {
	return write(w, http.StatusOK, response{Success: true})
}

func WriteJSON(w http.ResponseWriter, data any) error {
	return write(w, http.StatusOK, response{Success: true, Data: data})
}

func WriteError(w http.ResponseWriter, err error) error {
	// Default response for internal & unknown errors.
	statusCode := http.StatusInternalServerError
	response := response{Success: false, Data: http.StatusText(http.StatusInternalServerError)}

	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode != http.StatusInternalServerError {
		statusCode = apiErr.StatusCode
		response.Data = apiErr.Message
	}

	return write(w, statusCode, response)
}
