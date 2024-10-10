package todo

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoFilter struct {
	ID        *uuid.UUID `schema:"id"`
	UserID    *uuid.UUID `schema:"user_id"`
	Priority  *int       `schema:"priority"`
	DueDate   *time.Time `schema:"due_date"`
	Completed *bool      `schema:"completed"`
	Offset    int        `schema:"offset"`
	Limit     int        `schema:"limit"`
}

type TodoUpdate struct {
	Subject     *string    `json:"subject"`
	Description *string    `json:"description"`
	Priority    *int       `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
}
