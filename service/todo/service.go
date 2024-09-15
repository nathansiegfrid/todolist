package todo

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type repository interface {
	GetAll(ctx context.Context) ([]*Todo, error)
	Get(ctx context.Context, id uuid.UUID) (*Todo, error)
	Create(ctx context.Context, t *Todo) error
	Update(ctx context.Context, id uuid.UUID, t *Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	r repository
}

func NewService(r repository) *Service {
	return &Service{r}
}

func (s *Service) GetAllTodos(ctx context.Context) ([]*Todo, error) {
	return s.r.GetAll(ctx)
}

func (s *Service) GetTodoByID(ctx context.Context, id uuid.UUID) (*Todo, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) CreateTodo(ctx context.Context, req *Todo) error {
	// Validate user input.
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Description, validation.Required, validation.Length(0, 255)),
	); err != nil {
		return service.ErrValidation(err)
	}

	return s.r.Create(ctx, req)
}

func (s *Service) UpdateTodo(ctx context.Context, id uuid.UUID, req *Todo) error {
	// Validate user input.
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Description, validation.Required, validation.Length(0, 50)),
	); err != nil {
		return service.ErrValidation(err)
	}

	return s.r.Update(ctx, id, req)
}

func (s *Service) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}
