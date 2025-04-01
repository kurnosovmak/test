package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserWithEmailAlreadyExists    = errors.New("user with email already exists")
	ErrUserWithUsernameAlreadyExists = errors.New("user with username already exists")
)

type UserId = uuid.UUID

type User struct {
	ID           UserId    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUser struct {
	Username     string
	Email        string
	PasswordHash []byte
}

type UserRepository interface {
	Create(ctx context.Context, user *CreateUser) (UserId, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}
