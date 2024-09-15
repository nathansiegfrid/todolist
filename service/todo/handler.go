package todo

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) HTTPHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.getAllTodos)
	r.Get("/{id}", h.getTodoByID)
	r.Post("/", h.createTodo)
	r.Put("/{id}", h.updateTodo)
	r.Delete("/{id}", h.deleteTodo)
	return r
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.s.GetAllTodos(r.Context())
	if err != nil {
		service.WriteErr(w, err)
		return
	}

	service.WriteJSON(w, todos)
}

func (h *Handler) getTodoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteErr(w, service.ErrInvalidUUID(idStr))
		return
	}

	todo, err := h.s.GetTodoByID(r.Context(), id)
	if err != nil {
		service.WriteErr(w, err)
		return
	}

	service.WriteJSON(w, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req *Todo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		service.WriteErr(w, service.ErrInvalidJSON())
		return
	}

	if err := h.s.CreateTodo(r.Context(), req); err != nil {
		service.WriteErr(w, err)
		return
	}

	service.WriteOK(w)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteErr(w, service.ErrInvalidUUID(idStr))
		return
	}

	var req *Todo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		service.WriteErr(w, service.ErrInvalidJSON())
		return
	}

	if err := h.s.UpdateTodo(r.Context(), id, req); err != nil {
		service.WriteErr(w, err)
		return
	}

	service.WriteOK(w)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		service.WriteErr(w, service.ErrInvalidUUID(idStr))
		return
	}

	if err := h.s.DeleteTodo(r.Context(), id); err != nil {
		service.WriteErr(w, err)
		return
	}

	service.WriteOK(w)
}
