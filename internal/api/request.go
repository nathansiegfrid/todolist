package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
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
		return uuid.Nil, ErrInvalidID(idStr)
	}
	return id, nil
}

func ReadJSON[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return nil, ErrInvalidJSON()
	}
	return dst, nil
}

// ReadURLQuery maps URL query into struct using `schema` tags.
// Supports primitive types, time.Time, and uuid.UUID.
func ReadURLQuery[T any](r *http.Request) (*T, error) {
	dst := new(T)
	err := decodeURLQuery(r.URL.Query(), dst)
	if err != nil {
		return nil, ErrInvalidURLQuery(err)
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
