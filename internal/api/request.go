package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/samber/lo"
)

// ReadID reads {id} from URL path and parses it into uuid.UUID.
// Supports `go-chi/chi` and standard `http` routers.
func ReadID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, Errorf(http.StatusBadRequest, "Invalid ID param '%s'.", id)
	}
	return id, nil
}

func ReadJSON[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return nil, Error(http.StatusBadRequest, "Invalid JSON body.")
	}
	return dst, nil
}

// ReadURLQuery maps URL query into struct using `schema` tags.
// Supports primitive types, time.Time, and uuid.UUID.
func ReadURLQuery[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := schema.NewDecoder().Decode(dst, r.URL.Query())
	if err != nil {
		if errs, ok := err.(schema.MultiError); ok {
			// The MultiError map values doesn't make sense, so only the keys are returned.
			queryKeys := strings.Join(lo.Keys(errs), ", ")
			return nil, Errorf(http.StatusBadRequest, "Invalid URL query: %s.", queryKeys)
		}
		return nil, err // INTERNAL SERVER ERROR
	}
	return dst, nil
}
