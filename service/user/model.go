package user

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) SetNewPassword(p string) {
	sum := sha256.Sum256([]byte(p)) // Use checksum to avoid 72 bytes limit on bcrypt.
	h, err := bcrypt.GenerateFromPassword(sum[:], bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Errorf("error hashing password: %w", err))
	}
	u.PasswordHash = h
}

func (u *User) CheckPassword(p string) bool {
	sum := sha256.Sum256([]byte(p))
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, sum[:])
	return err == nil
}

type UserFilter struct {
	ID     *uuid.UUID `schema:"id"`
	Email  *string    `schema:"email"`
	Limit  int        `schema:"limit"`
	Offset int        `schema:"offset"`
}

// TODO: User update should require additional short-timed token, which requires re-auth with password.

type UserUpdate struct {
	Email    *string `json:"id"`
	Password *string `json:"password"`
}
