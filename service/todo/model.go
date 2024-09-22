package todo

import "github.com/google/uuid"

// TODO: Null UserID not allowed after user service is implemented.

type Todo struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.NullUUID `json:"user_id"`
	Description string        `json:"description"`
	Completed   bool          `json:"completed"`
}

type TodoFilter struct {
	UserID    *uuid.UUID `schema:"user_id"`
	Completed *bool      `schema:"completed"`
	Offset    int        `schema:"offset"`
	Limit     int        `schema:"limit"`
}

type TodoUpdate struct {
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
