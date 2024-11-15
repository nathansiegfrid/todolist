package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/handler"
	"github.com/nathansiegfrid/todolist/pkg/request"
	"github.com/nathansiegfrid/todolist/pkg/response"
	"github.com/nathansiegfrid/todolist/pkg/token"
)

var errLogin = response.Error(http.StatusUnauthorized, "Incorrect email or password.")

type repository interface {
	GetAll(ctx context.Context, filter *UserFilter) ([]*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, todo *User) error
	Update(ctx context.Context, id uuid.UUID, update *UserUpdate) error
}

type Handler struct {
	repository repository
	jwtAuth    *token.JWTAuth
}

func NewHandler(db *sql.DB, jwtAuth *token.JWTAuth) *Handler {
	return &Handler{
		repository: NewRepository(db),
		jwtAuth:    jwtAuth,
	}
}

func (h *Handler) HandleLoginRoute() http.HandlerFunc {
	return handler.MethodHandler{"POST": h.handleLogin()}.HandlerFunc()
}

func (h *Handler) HandleRegisterRoute() http.HandlerFunc {
	return handler.MethodHandler{"POST": h.handleRegister()}.HandlerFunc()
}

func (h *Handler) HandleVerifyAuthRoute() http.HandlerFunc {
	return handler.MethodHandler{"GET": h.handleVerifyAuth()}.HandlerFunc()
}

func (h *Handler) handleLogin() http.HandlerFunc {
	type requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseData struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	return handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		// Read request body.
		reqBody, err := request.ReadJSON[requestData](r)
		if err != nil {
			return err
		}

		users, err := h.repository.GetAll(r.Context(), &UserFilter{Email: &reqBody.Email, Limit: 1})
		if err != nil {
			return err
		}

		if len(users) == 0 || !users[0].CheckPassword(reqBody.Password) {
			return errLogin
		}

		token, err := h.jwtAuth.GenerateToken(users[0].ID, 5*time.Minute)
		if err != nil {
			return err
		}
		refreshToken, err := h.jwtAuth.GenerateToken(users[0].ID, 72*time.Hour)
		if err != nil {
			return err
		}

		return response.WriteJSON(w, &responseData{
			Token:        token,
			RefreshToken: refreshToken,
		})
	})
}

func (h *Handler) handleRegister() http.HandlerFunc {
	type requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		// Read request body.
		reqBody, err := request.ReadJSON[requestData](r)
		if err != nil {
			return err
		}

		// Validate user input.
		if err := validation.ValidateStruct(reqBody,
			validation.Field(&reqBody.Email, validation.Required, is.Email),
			validation.Field(&reqBody.Password, validation.Required, validation.Length(8, 0)),
		); err != nil {
			if errs, ok := err.(validation.Errors); ok {
				return response.ErrDataValidation(errs)
			}
			return err
		}

		// Create user entity from request.
		user := &User{Email: reqBody.Email}
		user.SetNewPassword(reqBody.Password)

		err = h.repository.Create(r.Context(), user)
		if err != nil {
			return err
		}

		return response.WriteOK(w)
	})
}

// HandleVerifyAuth returns user info if the request is correctly authenticated.
// Use with Authenticator middleware.
func (h *Handler) handleVerifyAuth() http.HandlerFunc {
	type responseData struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	return handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		userID := request.UserIDFromContext(r.Context())
		user, err := h.repository.Get(r.Context(), userID)
		if err != nil {
			return err
		}

		return response.WriteJSON(w, &responseData{
			ID:    user.ID.String(),
			Email: user.Email,
		})
	})
}
