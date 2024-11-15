package handler

import (
	"fmt"
	"net/http"

	"github.com/nathansiegfrid/todolist/pkg/logger"
	"github.com/nathansiegfrid/todolist/pkg/response"
)

// ErrorHandlerFunc wraps a handler function that returns an error into an http.HandlerFunc.
func ErrorHandlerFunc(serveWithError func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := serveWithError(w, r)
		if err == nil {
			return
		}

		ctx := r.Context()
		res := response.ErrorResponseFrom(err)
		if res.StatusCode >= http.StatusInternalServerError {
			logger.Error(ctx, fmt.Sprintf("Unexpected error: %s.", err), "category", "internal_error")
		} else {
			logger.Info(ctx, res.Message, "category", "client_error")
		}
		response.WriteError(w, res)
	}
}
