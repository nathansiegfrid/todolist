package api

import (
	"encoding/json"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Optional can distinguish between "null" value and undefined key.
type Optional[T any] struct {
	Value   T    `schema:"value"`
	Defined bool `schema:"-"`
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (t *Optional[T]) UnmarshalJSON(value []byte) error {
	t.Defined = true
	return json.Unmarshal(value, &t.Value)
}

// UnmarshalText implements encoding.TextUnmarshaler interface to support `gorilla/schema` decoding.
func (t *Optional[T]) UnmarshalText(value []byte) error {
	t.Defined = true
	valueStr := string(value)
	// If value is "null" and T is a pointer, leave it as nil.
	if valueStr == "null" && reflect.TypeOf(t.Value).Kind() == reflect.Ptr {
		return nil
	}
	// Create artificial map to for schema decoder.
	src := map[string][]string{"value": {valueStr}}
	return decodeURLQuery(src, t)
}

// // Value implements the `driver.Valuer` interface to support `ozzo-validation`.
// func (t Optional[T]) Value() (driver.Value, error) {
// 	// Don't change the receiver to pointer because it won't work with `ozzo-validation`.
// 	if !t.Defined {
// 		return nil, nil
// 	}
// 	return t.T, nil
// }

// OptionalValidator wraps `ozzo-validation` rules and applies them to Optional value.
type OptionalValidator[T any] struct {
	Rules []validation.Rule
}

func NewOptionalValidator[T any](rules ...validation.Rule) validation.Rule {
	return OptionalValidator[T]{Rules: rules}
}

func (v OptionalValidator[T]) Validate(value any) error {
	t, ok := value.(Optional[T])
	if ok {
		// Skip validation for undefined JSON key.
		if !t.Defined {
			return nil
		}
		value = t.Value
	}

	for _, rule := range v.Rules {
		if err := rule.Validate(value); err != nil {
			return err
		}
	}
	return nil
}
