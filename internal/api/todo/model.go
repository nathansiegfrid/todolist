package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/internal/api"
)

type Todo struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TodoUpdate struct {
	Subject     api.Optional[string]     `json:"subject"`
	Description api.Optional[string]     `json:"description"`
	Priority    api.Optional[int]        `json:"priority"`
	DueDate     api.Optional[*time.Time] `json:"due_date"`
	Completed   api.Optional[bool]       `json:"completed"`
}

type TodoFilter struct {
	ID        api.Optional[uuid.UUID]  `schema:"id"`
	UserID    api.Optional[uuid.UUID]  `schema:"user_id"`
	Priority  api.Optional[int]        `schema:"priority"`
	DueDate   api.Optional[*time.Time] `schema:"due_date"`
	Completed api.Optional[bool]       `schema:"completed"`
	Offset    int                      `schema:"offset"`
	Limit     int                      `schema:"limit"`
}
