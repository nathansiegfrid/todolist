package response

import (
	"encoding/json"
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

func write(w http.ResponseWriter, statusCode int, body responseBody) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(&body)
}

func WriteOK(w http.ResponseWriter) error {
	return write(w, http.StatusOK, responseBody{Status: "SUCCESS"})
}

func WriteJSON(w http.ResponseWriter, data any) error {
	return write(w, http.StatusOK, responseBody{Status: "SUCCESS", Data: data})
}

func WriteError(w http.ResponseWriter, res ErrorResponse) error {
	var status string
	if res.StatusCode >= http.StatusInternalServerError {
		status = "ERROR"
	} else {
		status = "FAIL"
	}
	return write(w, res.StatusCode, responseBody{
		Status:  status,
		Message: res.Message,
		Data:    res.Data,
	})
}
