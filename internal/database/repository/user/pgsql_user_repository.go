package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kurnosovmak/test/internal/auth"
	"github.com/kurnosovmak/test/internal/database"
)

type PgsqlUserRepository struct {
	db *database.Database
}

func NewPgsqlUserRepository(db *database.Database) *PgsqlUserRepository {
	return &PgsqlUserRepository{db: db}
}
func (r *PgsqlUserRepository) Create(ctx context.Context, user *auth.CreateUser) (auth.UserId, error) {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return auth.UserId{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	row := tx.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash)
	var id auth.UserId
	if err := row.Scan(&id); err != nil {
		if database.IsDuplicatedKeyError(err, "users_email") {
			return auth.UserId{}, auth.ErrUserWithEmailAlreadyExists
		}
		if database.IsDuplicatedKeyError(err, "users_username") {
			return auth.UserId{}, auth.ErrUserWithUsernameAlreadyExists
		}

		return auth.UserId{}, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *PgsqlUserRepository) GetByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `
		SELECT 
		id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1 
		LIMIT 1
	`

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	row := tx.QueryRow(ctx, query, email)
	var user auth.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}
