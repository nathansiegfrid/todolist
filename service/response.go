package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

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
	statusCode := http.StatusOK
	response := response{Success: true}

	return write(w, statusCode, response)
}

func WriteJSON(w http.ResponseWriter, data any) error {
	statusCode := http.StatusOK
	response := response{Success: true, Data: data}

	return write(w, statusCode, response)
}

func WriteErr(w http.ResponseWriter, err error) error {
	// Default response for internal & unknown errors.
	statusCode := http.StatusInternalServerError
	response := response{Success: false, Data: "internal server error"}

	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode != http.StatusInternalServerError {
		statusCode = apiErr.StatusCode
		response.Data = apiErr.Message
	} else {
		log.Print(err)
	}

	return write(w, statusCode, response)
}
