package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
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
	Subject     service.Optional[string]     `json:"subject"`
	Description service.Optional[string]     `json:"description"`
	Priority    service.Optional[int]        `json:"priority"`
	DueDate     service.Optional[*time.Time] `json:"due_date"`
	Completed   service.Optional[bool]       `json:"completed"`
}

type TodoFilter struct {
	ID        service.Optional[uuid.UUID]  `schema:"id"`
	UserID    service.Optional[uuid.UUID]  `schema:"user_id"`
	Priority  service.Optional[int]        `schema:"priority"`
	DueDate   service.Optional[*time.Time] `schema:"due_date"`
	Completed service.Optional[bool]       `schema:"completed"`
	Offset    int                          `schema:"offset"`
	Limit     int                          `schema:"limit"`
}
