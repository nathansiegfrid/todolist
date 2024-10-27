package api

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// Optional can distinguish between "null" value and undefined key.
type Optional[T any] struct {
	Data    T    `schema:"data"`
	Defined bool `schema:"-"`
}

func NewOptional[T any](data T) Optional[T] {
	return Optional[T]{Data: data, Defined: true}
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (t *Optional[T]) UnmarshalJSON(data []byte) error {
	t.Defined = true
	return json.Unmarshal(data, &t.Data)
}

// UnmarshalText implements encoding.TextUnmarshaler interface to support `gorilla/schema` decoding.
func (t *Optional[T]) UnmarshalText(data []byte) error {
	t.Defined = true
	// If data is "null" and T is a pointer, leave it as nil.
	if string(data) == "null" && reflect.TypeOf(t.Data).Kind() == reflect.Ptr {
		return nil
	}
	// Create artificial map to use schema decoder.
	src := map[string][]string{"data": {string(data)}}
	return decodeURLQuery(src, t)
}

// Value implements the `driver.Valuer` interface to support `ozzo-validation`.
func (t Optional[T]) Value() (driver.Value, error) {
	// Don't use pointer receiver because it won't work with `ozzo-validation`.
	if !t.Defined {
		return nil, nil
	}
	return t.Data, nil
}
