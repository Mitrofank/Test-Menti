package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	signingKey string
	ttl        time.Duration
}

func NewTokenManager(signingKey string, ttl time.Duration) (*TokenManager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}
	return &TokenManager{signingKey: signingKey, ttl: ttl}, nil
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

func (tm *TokenManager) GenerateToken(userID int, roleID int) (string, error) {
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		RoleID: roleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tm.signingKey))
}
