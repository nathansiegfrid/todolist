package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/nathansiegfrid/todolist/pkg/field"
)

// Types in `guregu/null` package implements `json.Unmarshaler` and `encoding.TextUnmarshaler` interfaces.
// They supports URL query parsing with `gorilla/schema` decoder.

type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	DueDate     null.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoUpdate struct {
	Subject     field.Option[string]    `json:"subject"`
	Description field.Option[string]    `json:"description"`
	Priority    field.Option[int]       `json:"priority"`
	DueDate     field.Option[null.Time] `json:"due_date"`
	Completed   field.Option[bool]      `json:"completed"`
}

type TodoFilter struct {
	ID        *uuid.UUID `schema:"id"`
	UserID    *uuid.UUID `schema:"user_id"`
	Priority  *int       `schema:"priority"`
	DueDate   *null.Time `schema:"due_date"`
	Completed *bool      `schema:"completed"`
	Offset    int        `schema:"offset"`
	Limit     int        `schema:"limit"`
}
