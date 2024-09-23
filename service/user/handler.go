package user

import (
	"context"
	"database/sql"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type repository interface {
	GetAll(ctx context.Context, filter *UserFilter) ([]*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, todo *User) error
	Update(ctx context.Context, id uuid.UUID, update *UserUpdate) error
}

type Handler struct {
	repository repository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repository: NewRepository(db),
	}
}

func (h *Handler) HandleLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body.
		reqBody, err := service.ReadJSON[request](r)
		if err != nil {
			service.LogInfo(r.Context(), err)
			service.WriteError(w, service.ErrInvalidJSON())
			return
		}

		users, err := h.repository.GetAll(r.Context(), &UserFilter{Email: &reqBody.Email, Limit: 1})
		if err != nil {
			service.LogInternalError(r.Context(), err)
			service.WriteError(w, err)
			return
		}

		if len(users) == 0 || !users[0].CheckPassword(reqBody.Password) {
			service.WriteError(w, service.Error(http.StatusUnauthorized, "incorrect email or password"))
			return
		}

		service.WriteJSON(w, &response{
			Token:        "not_implemented",
			RefreshToken: "not_implemented",
		})
	}
}

func (h *Handler) HandleRegister() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body.
		reqBody, err := service.ReadJSON[request](r)
		if err != nil {
			service.LogInfo(r.Context(), err)
			service.WriteError(w, service.ErrInvalidJSON())
			return
		}

		// Validate user input.
		err = validation.ValidateStruct(reqBody,
			validation.Field(&reqBody.Email, validation.Required, is.Email),
			validation.Field(&reqBody.Password, validation.Required, validation.Length(8, 0)),
		)
		if err != nil {
			if valErr, ok := err.(validation.Errors); ok {
				service.WriteError(w, service.ErrValidation(valErr))
			} else {
				service.LogError(r.Context(), err)
				service.WriteError(w, err)
			}
			return
		}

		// Create user entity from request.
		user := &User{Email: reqBody.Email}
		user.SetNewPassword(reqBody.Password)

		err = h.repository.Create(r.Context(), user)
		if err != nil {
			service.LogInternalError(r.Context(), err)
			service.WriteError(w, err)
			return
		}
		service.WriteOK(w)
	}
}
