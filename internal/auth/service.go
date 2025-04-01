package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Service interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (UserId, error)
}

type authService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) Service {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) RegisterUser(ctx context.Context, input RegisterUserInput) (UserId, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserId{}, err
	}

	// Создаем нового пользователя
	createUser := &CreateUser{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.Create(ctx, createUser)
}
