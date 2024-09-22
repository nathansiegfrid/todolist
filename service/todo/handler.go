package todo

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
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

func (h *Handler) HTTPHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.getAllTodos)
	r.Get("/{id}", h.getTodo)
	r.Post("/", h.createTodo)
	r.Patch("/{id}", h.updateTodo)
	r.Delete("/{id}", h.deleteTodo)
	return r
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	// Read URL query.
	filter, err := service.ReadURLQuery[TodoFilter](r)
	if err != nil {
		service.LogError(r.Context(), err)
	}

	todos, err := h.repository.GetAll(r.Context(), filter)
	if err != nil {
		service.LogInternalError(r.Context(), err)
		service.WriteError(w, err)
		return
	}
	service.WriteJSON(w, todos)
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param.
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteError(w, service.ErrInvalidID(idStr))
		return
	}

	todo, err := h.repository.Get(r.Context(), id)
	if err != nil {
		service.LogInternalError(r.Context(), err)
		service.WriteError(w, err)
		return
	}
	service.WriteJSON(w, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	todo, err := service.ReadJSON[Todo](r)
	if err != nil {
		service.LogError(r.Context(), err)
		service.WriteError(w, service.ErrInvalidJSON())
		return
	}

	// Validate user input.
	err = validation.ValidateStruct(todo,
		validation.Field(&todo.Description, validation.Required, validation.Length(0, 255)),
	)
	if err != nil {
		service.WriteError(w, service.ErrValidation(err))
		return
	}

	err = h.repository.Create(r.Context(), todo)
	if err != nil {
		service.LogInternalError(r.Context(), err)
		service.WriteError(w, err)
		return
	}
	service.WriteOK(w)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param.
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteError(w, service.ErrInvalidID(idStr))
		return
	}

	// Read request body.
	update, err := service.ReadJSON[TodoUpdate](r)
	if err != nil {
		service.LogError(r.Context(), err)
		service.WriteError(w, service.ErrInvalidJSON())
		return
	}

	// Validate user input.
	err = validation.ValidateStruct(update,
		validation.Field(&update.Description, validation.Required, validation.Length(0, 255)),
	)
	if err != nil {
		service.WriteError(w, service.ErrValidation(err))
		return
	}

	err = h.repository.Update(r.Context(), id, update)
	if err != nil {
		service.LogInternalError(r.Context(), err)
		service.WriteError(w, err)
		return
	}
	service.WriteOK(w)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Read request param.
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteError(w, service.ErrInvalidID(idStr))
		return
	}

	err = h.repository.Delete(r.Context(), id)
	if err != nil {
		service.LogInternalError(r.Context(), err)
		service.WriteError(w, err)
		return
	}
	service.WriteOK(w)
}
