package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// responseBody standardizes the response format.
// Status is either "SUCCESS", "FAIL", or "ERROR".
// Status "FAIL" is used for client errors.
// Status "ERROR" is used for server errors.
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
	var res ErrorResponse
	if errors.As(err, &res) && res.StatusCode != http.StatusInternalServerError {
		return write(w, res.StatusCode, &responseBody{
			Status:  "FAIL",
			Message: res.Message,
			Data:    res.Data,
		})
	}

	return write(w, http.StatusInternalServerError, &responseBody{
		Status:  "ERROR",
		Message: "Unexpected error. We've noted the issue. Please try again later.", // Log the error.
	})
}
