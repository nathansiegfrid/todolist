package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/internal/api"
)

type repository interface {
	GetAll(ctx context.Context, filter *UserFilter) ([]*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, todo *User) error
	Update(ctx context.Context, id uuid.UUID, update *UserUpdate) error
}

type tokenGenerator interface {
	GenerateToken(subject string, duration time.Duration) (signedToken string, err error)
}

type Handler struct {
	repository     repository
	tokenGenerator tokenGenerator
}

func NewHandler(db *sql.DB, tokenGenerator tokenGenerator) *Handler {
	return &Handler{
		repository:     NewRepository(db),
		tokenGenerator: tokenGenerator,
	}
}

func (h *Handler) HandleLoginRoute() http.HandlerFunc {
	return api.MethodHandler{"POST": h.handleLogin()}.HandlerFunc()
}

func (h *Handler) HandleRegisterRoute() http.HandlerFunc {
	return api.MethodHandler{"POST": h.handleRegister()}.HandlerFunc()
}

func (h *Handler) HandleVerifyAuthRoute() http.HandlerFunc {
	return api.MethodHandler{"GET": h.handleVerifyAuth()}.HandlerFunc()
}

func (h *Handler) handleLogin() http.HandlerFunc {
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
		reqBody, err := api.ReadJSON[request](r)
		if err != nil {
			api.LogInfo(r.Context(), err)
			api.WriteError(w, api.ErrInvalidJSON())
			return
		}

		users, err := h.repository.GetAll(r.Context(), &UserFilter{Email: &reqBody.Email, Limit: 1})
		if err != nil {
			api.LogErrorInternal(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		if len(users) == 0 || !users[0].CheckPassword(reqBody.Password) {
			err := api.Error(http.StatusUnauthorized, "Incorrect email or password.")
			api.WriteError(w, err)
			return
		}

		userID := users[0].ID.String()
		token, err := h.tokenGenerator.GenerateToken(userID, 5*time.Minute)
		if err != nil {
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}
		refreshToken, err := h.tokenGenerator.GenerateToken(userID, 72*time.Hour)
		if err != nil {
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		api.WriteJSON(w, &response{
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

func (h *Handler) handleRegister() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body.
		reqBody, err := api.ReadJSON[request](r)
		if err != nil {
			api.LogInfo(r.Context(), err)
			api.WriteError(w, api.ErrInvalidJSON())
			return
		}

		// Validate user input.
		err = validation.ValidateStruct(reqBody,
			validation.Field(&reqBody.Email, validation.Required, is.Email),
			validation.Field(&reqBody.Password, validation.Required, validation.Length(8, 0)),
		)
		if err != nil {
			if valErr, ok := err.(validation.Errors); ok {
				api.WriteError(w, api.ErrValidation(valErr))
			} else {
				api.LogError(r.Context(), err)
				api.WriteError(w, err)
			}
			return
		}

		// Create user entity from request.
		user := &User{Email: reqBody.Email}
		user.SetNewPassword(reqBody.Password)

		err = h.repository.Create(r.Context(), user)
		if err != nil {
			api.LogErrorInternal(r.Context(), err)
			api.WriteError(w, err)
			return
		}
		api.WriteOK(w)
	}
}

// HandleVerifyAuth returns user info if the request is correctly authenticated.
// Use with Authenticator middleware.
func (h *Handler) handleVerifyAuth() http.HandlerFunc {
	type response struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := api.UserIDFromContext(r.Context())
		user, err := h.repository.Get(r.Context(), userID)
		if err != nil {
			api.LogErrorInternal(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		api.WriteJSON(w, &response{
			ID:    user.ID.String(),
			Email: user.Email,
		})
	}
}
