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

var errLogin = api.Error(http.StatusUnauthorized, "Incorrect email or password.")

type repository interface {
	GetAll(ctx context.Context, filter *UserFilter) ([]*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, todo *User) error
	Update(ctx context.Context, id uuid.UUID, update *UserUpdate) error
}

type Handler struct {
	repository repository
	jwtService *JWTService
}

func NewHandler(db *sql.DB, jwtService *JWTService) *Handler {
	return &Handler{
		repository: NewRepository(db),
		jwtService: jwtService,
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
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		users, err := h.repository.GetAll(r.Context(), &UserFilter{Email: &reqBody.Email, Limit: 1})
		if err != nil {
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		if len(users) == 0 || !users[0].CheckPassword(reqBody.Password) {
			err := errLogin
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		token, err := h.jwtService.GenerateToken(users[0].ID, 5*time.Minute)
		if err != nil {
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}
		refreshToken, err := h.jwtService.GenerateToken(users[0].ID, 72*time.Hour)
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
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		// Validate user input.
		if err := validation.ValidateStruct(reqBody,
			validation.Field(&reqBody.Email, validation.Required, is.Email),
			validation.Field(&reqBody.Password, validation.Required, validation.Length(8, 0)),
		); err != nil {
			err := api.ErrDataValidation(err)
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		// Create user entity from request.
		user := &User{Email: reqBody.Email}
		user.SetNewPassword(reqBody.Password)

		err = h.repository.Create(r.Context(), user)
		if err != nil {
			api.LogError(r.Context(), err)
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
			api.LogError(r.Context(), err)
			api.WriteError(w, err)
			return
		}

		api.WriteJSON(w, &response{
			ID:    user.ID.String(),
			Email: user.Email,
		})
	}
}
