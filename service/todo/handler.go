package todo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type repository interface {
	GetAll(ctx context.Context) ([]*Todo, error)
	Get(ctx context.Context, id uuid.UUID) (*Todo, error)
	Create(ctx context.Context, req *CreateTodoRequest) error
	Update(ctx context.Context, id uuid.UUID, req *UpdateTodoRequest) error
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
	r.Put("/{id}", h.updateTodo)
	r.Delete("/{id}", h.deleteTodo)
	return r
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repository.GetAll(r.Context())
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
		service.WriteError(w, service.ErrInvalidUUID(idStr))
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
	var reqBody *CreateTodoRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		service.WriteError(w, service.ErrInvalidJSON())
		return
	}

	// Validate user input.
	err = validation.ValidateStruct(reqBody,
		validation.Field(&reqBody.Description, validation.Required, validation.Length(0, 255)),
	)
	if err != nil {
		service.WriteError(w, service.ErrValidation(err))
		return
	}

	err = h.repository.Create(r.Context(), reqBody)
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
		service.WriteError(w, service.ErrInvalidUUID(idStr))
		return
	}

	// Read request body.
	var reqBody *UpdateTodoRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		service.WriteError(w, service.ErrInvalidJSON())
		return
	}

	// Validate user input.
	err = validation.ValidateStruct(reqBody,
		validation.Field(&reqBody.Description, validation.Required, validation.Length(0, 255)),
	)
	if err != nil {
		service.WriteError(w, service.ErrValidation(err))
		return
	}

	err = h.repository.Update(r.Context(), id, reqBody)
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
		service.WriteError(w, service.ErrInvalidUUID(idStr))
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
