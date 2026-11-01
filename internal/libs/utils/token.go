package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func SignToken(userID, tokenType, tokenID, secret string, lifetime time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(lifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseToken(tokenStr, tokenType, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.Subject == "" {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.Type != tokenType {
		return nil, errors.New("unexpected token type")
	}
	return claims, nil
}

func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func VerificationTokenExpiry() time.Time {
	return time.Now().UTC().Add(time.Hour)
}

func PasswordResetTokenExpiry() time.Time {
	return time.Now().UTC().Add(time.Hour)
}
