package api

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Optional can distinguish between "null" value and undefined key.
type Optional[T any] struct {
	Value   T
	Defined bool
}

func (t *Optional[T]) UnmarshalJSON(value []byte) error {
	t.Defined = true
	return json.Unmarshal(value, &t.Value)
}

// UnmarshalText implements encoding.TextUnmarshaler interface to support `gorilla/schema`.
func (t *Optional[T]) UnmarshalText(value []byte) error {
	t.Defined = true
	return json.Unmarshal(value, &t.Value)
}

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
