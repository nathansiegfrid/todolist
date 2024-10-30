package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

// ReadID reads {id} from URL path and parses it into uuid.UUID.
// Its centralized here to allow easy modifications if the routing library changes.
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
	err := decodeURLQuery(r.URL.Query(), dst)
	if err != nil {
		var errs schema.MultiError
		if errors.As(err, &errs) {
			// The multi error value doesn't make sense, so only the keys are returned.
			keys := make([]string, 0, len(errs))
			for k := range errs {
				keys = append(keys, k)
			}
			return nil, Errorf(http.StatusBadRequest, "Invalid URL query: %s.", strings.Join(keys, ", "))
		}
		return nil, err // Internal server error.
	}
	return dst, nil
}

func decodeURLQuery(src url.Values, dst any) error {
	dec := schema.NewDecoder()
	// Add custom converter because default decoder doesn't support time.Time.
	dec.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		t, err := time.Parse(time.RFC3339, value) // Same time format as JSON.
		if err != nil {
			return reflect.Value{} // Zero value represents invalid Value and "Invalid" kind.
		}
		return reflect.ValueOf(t)
	})
	return dec.Decode(dst, src)
}
