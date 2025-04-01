package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/kurnosovmak/test/pkg/jwt"
	"github.com/kurnosovmak/test/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}

type Service interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (UserId, error)
	LoginUser(ctx context.Context, input LoginUserInput) (*LoginResponse, error)
	ValidateToken(token string) (*jwt.UserClaims, error)
}

type authService struct {
	userRepo   UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo UserRepository, jwtManager *jwt.JWTManager) Service {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
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

func (s *authService) LoginUser(ctx context.Context, input LoginUserInput) (*LoginResponse, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Возвращаем ошибку, если пользователь не найден
			return nil, ErrInvalidCredentials
		}
		logger.Error("failed to get user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.jwtManager.GenerateToken(jwt.Payload{
		UserID:   user.ID.String(),
		Username: user.Username,
	})
	if err != nil {
		logger.Error("failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

func (s *authService) ValidateToken(token string) (*jwt.UserClaims, error) {
	return s.jwtManager.ValidateToken(token)
}
