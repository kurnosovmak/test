package jwt

import (
	"context"
	"net/http"
	"strings"

	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDContextKey contextKey = "user_id"
	AuthHeaderKey    string     = "Authorization"
	BearerSchema     string     = "Bearer "
)

// AuthMiddleware проверяет JWT токен и добавляет ID пользователя в контекст
func AuthMiddleware(next http.Handler, jwtService *JWTManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Получаем токен из заголовка
		authHeader := r.Header.Get(AuthHeaderKey)
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing authorization header"))
			return
		}

		// Проверяем формат токена
		if !strings.HasPrefix(authHeader, BearerSchema) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token format"))
			return
		}

		// Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		// Проверяем токен
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("expired token"))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}

		// Добавляем ID пользователя в контекст
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		return uuid.UUID{}, errors.New("user ID not found in context")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid user ID")
	}
	return uid, nil
}
