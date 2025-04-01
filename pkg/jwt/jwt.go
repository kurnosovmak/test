package jwt

import (
	"errors"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwtLib.RegisteredClaims
	Payload
}

type Payload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (m *JWTManager) GenerateToken(payload Payload) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
		Payload: payload,
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(tokenStr string) (*UserClaims, error) {
	token, err := jwtLib.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwtLib.Token) (interface{}, error) {
			_, ok := token.Method.(*jwtLib.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidToken
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwtLib.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
