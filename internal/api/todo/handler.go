package todo

import (
	"context"
	"database/sql"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/internal/api"
)

type repository interface {
	GetAll(ctx context.Context, filter *TodoFilter) ([]*Todo, error)
	Get(ctx context.Context, id uuid.UUID) (*Todo, error)
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, id uuid.UUID, update *TodoUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	repository repository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repository: NewRepository(db),
	}
}

func (h *Handler) HandleTodosRoute() http.HandlerFunc {
	return api.MethodHandler{
		"GET":  h.getAllTodos,
		"POST": h.createTodo,
	}.HandlerFunc()
}

func (h *Handler) HandleTodosIDRoute() http.HandlerFunc {
	return api.MethodHandler{
		"GET":    h.getTodo,
		"PATCH":  h.updateTodo,
		"DELETE": h.deleteTodo,
	}.HandlerFunc()
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	// Read URL query.
	filter, err := api.ReadURLQuery[TodoFilter](r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	todos, err := h.repository.GetAll(r.Context(), filter)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}
	api.WriteJSON(w, todos)
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param "id".
	id, err := api.ReadID(r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	todo, err := h.repository.Get(r.Context(), id)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}
	api.WriteJSON(w, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	todo, err := api.ReadJSON[Todo](r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	// Validate user input.
	if err := validation.ValidateStruct(todo,
		validation.Field(&todo.Subject, validation.Required, validation.Length(0, 100)),
		validation.Field(&todo.Description, validation.Length(0, 1000)),
	); err != nil {
		err := api.ErrValidation(err)
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	err = h.repository.Create(r.Context(), todo)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}
	api.WriteOK(w)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param "id".
	id, err := api.ReadID(r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	// Read request body.
	update, err := api.ReadJSON[TodoUpdate](r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	// Validate user input.
	if err := validation.ValidateStruct(update,
		validation.Field(&update.Subject, validation.NilOrNotEmpty, validation.Length(0, 100)),
		validation.Field(&update.Description, validation.Length(0, 1000)),
	); err != nil {
		err := api.ErrValidation(err)
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	err = h.repository.Update(r.Context(), id, update)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}
	api.WriteOK(w)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param "id".
	id, err := api.ReadID(r)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}

	err = h.repository.Delete(r.Context(), id)
	if err != nil {
		api.LogError(r.Context(), err)
		api.WriteError(w, err)
		return
	}
	api.WriteOK(w)
}
