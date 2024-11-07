package field

import (
	"database/sql/driver"
	"encoding/json"
)

// Option is used to represent optional fields in JSON.
// It can distinguish between missing fields and fields with "null" values.
type Option[T any] struct {
	value   T
	defined bool
}

func OptionFrom[T any](data T) Option[T] {
	return Option[T]{data, true}
}

// Value implements the `driver.Valuer` interface.
func (t Option[T]) Value() (driver.Value, error) {
	// This method is implemented to support `ozzo-validation`.
	// Validation rules will be applied to the returned value.
	if !t.defined {
		return nil, nil
	}
	return t.value, nil
}

func (t Option[T]) ValueOr(fallback T) T {
	if !t.defined {
		return fallback
	}
	return t.value
}

func (t Option[T]) ValueOrZero() T {
	if !t.defined {
		var zero T
		return zero
	}
	return t.value
}

func (t Option[T]) Defined() bool {
	return t.defined
}

// MarshalJSON implements the `json.Marshaler` interface.
func (t Option[T]) MarshalJSON() ([]byte, error) {
	if !t.defined {
		return []byte("null"), nil
	}
	return json.Marshal(t.value)
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (t *Option[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &t.value); err != nil {
		return err
	}
	t.defined = true
	return nil
}
