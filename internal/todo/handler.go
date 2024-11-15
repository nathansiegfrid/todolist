package todo

import (
	"context"
	"database/sql"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/handler"
	"github.com/nathansiegfrid/todolist/pkg/request"
	"github.com/nathansiegfrid/todolist/pkg/response"
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
	return handler.MethodHandler{
		"GET":  handler.ErrorHandlerFunc(h.getAllTodos),
		"POST": handler.ErrorHandlerFunc(h.createTodo),
	}.HandlerFunc()
}

func (h *Handler) HandleTodosIDRoute() http.HandlerFunc {
	return handler.MethodHandler{
		"GET":    handler.ErrorHandlerFunc(h.getTodo),
		"PATCH":  handler.ErrorHandlerFunc(h.updateTodo),
		"DELETE": handler.ErrorHandlerFunc(h.deleteTodo),
	}.HandlerFunc()
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) error {
	// Read URL query.
	filter, err := request.ReadURLQuery[TodoFilter](r)
	if err != nil {
		return err
	}

	todos, err := h.repository.GetAll(r.Context(), filter)
	if err != nil {
		return err
	}

	return response.WriteJSON(w, todos)
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) error {
	// Read request param "id".
	id, err := request.ReadID(r)
	if err != nil {
		return err
	}

	todo, err := h.repository.Get(r.Context(), id)
	if err != nil {
		return err
	}

	return response.WriteJSON(w, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) error {
	// Read request body.
	todo, err := request.ReadJSON[Todo](r)
	if err != nil {
		return err
	}

	// Validate user input.
	if err := validation.ValidateStruct(todo,
		validation.Field(&todo.Subject, validation.Required, validation.Length(0, 100)),
		validation.Field(&todo.Description, validation.Length(0, 1000)),
	); err != nil {
		if errs, ok := err.(validation.Errors); ok {
			return response.ErrDataValidation(errs)
		}
		return err
	}

	err = h.repository.Create(r.Context(), todo)
	if err != nil {
		return err
	}

	return response.WriteOK(w)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) error {
	// Read request param "id".
	id, err := request.ReadID(r)
	if err != nil {
		return err
	}

	// Read request body.
	update, err := request.ReadJSON[TodoUpdate](r)
	if err != nil {
		return err
	}

	// Validate user input.
	if err := validation.ValidateStruct(update,
		validation.Field(&update.Subject, validation.NilOrNotEmpty, validation.Length(0, 100)),
		validation.Field(&update.Description, validation.Length(0, 1000)),
	); err != nil {
		if errs, ok := err.(validation.Errors); ok {
			return response.ErrDataValidation(errs)
		}
		return err
	}

	err = h.repository.Update(r.Context(), id, update)
	if err != nil {
		return err
	}

	return response.WriteOK(w)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) error {
	// Read request param "id".
	id, err := request.ReadID(r)
	if err != nil {
		return err
	}

	err = h.repository.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	return response.WriteOK(w)
}
