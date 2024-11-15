package request

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/nathansiegfrid/todolist/pkg/response"
	"github.com/samber/lo"
)

// ReadID reads {id} from URL path and parses it into uuid.UUID.
// Supports `go-chi/chi` and standard `http` routers.
func ReadID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, response.Errorf(http.StatusBadRequest, "Invalid ID param '%s'.", id)
	}
	return id, nil
}

func ReadJSON[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return nil, response.Error(http.StatusBadRequest, "Invalid JSON body.")
	}
	return dst, nil
}

// ReadURLQuery maps URL query into struct using `schema` tags.
// Supports primitive types, time.Time, and uuid.UUID.
func ReadURLQuery[T any](r *http.Request) (*T, error) {
	dst := new(T)
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err := dec.Decode(dst, r.URL.Query())
	if err != nil {
		if errs, ok := err.(schema.MultiError); ok {
			// The MultiError map values doesn't make sense, so only the keys are returned.
			return nil, response.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid URL query.",
				Data:       lo.Keys(errs),
			}
		}
		return nil, err // INTERNAL SERVER ERROR
	}
	return dst, nil
}
