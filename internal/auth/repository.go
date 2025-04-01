package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrUserWithEmailAlreadyExists    = errors.New("user with email already exists")
	ErrUserWithUsernameAlreadyExists = errors.New("user with username already exists")
)

type UserId = uuid.UUID

type User struct {
	ID           UserId
	Username     string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUser struct {
	Username     string
	Email        string
	PasswordHash []byte
}

type UserRepository interface {
	Create(ctx context.Context, user *CreateUser) (UserId, error)
	// GetByID(ctx context.Context, id UseerId) (*User, error)
	// GetByUsername(ctx context.Context, username string) (*User, error)
	// GetByEmail(ctx context.Context, email string) (*User, error)
	// Update(ctx context.Context, user *User) error
}
