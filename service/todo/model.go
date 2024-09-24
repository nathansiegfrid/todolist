package todo

import "github.com/google/uuid"

type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
}

type TodoFilter struct {
	ID        *uuid.UUID `schema:"id"`
	UserID    *uuid.UUID `schema:"user_id"`
	Completed *bool      `schema:"completed"`
	Offset    int        `schema:"offset"`
	Limit     int        `schema:"limit"`
}

type TodoUpdate struct {
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
