package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTAuth struct {
	secret   string
	duration time.Duration
}

func NewJWTAuth(secret string, duration time.Duration) *JWTAuth {
	return &JWTAuth{
		secret:   secret,
		duration: duration,
	}
}

func (a *JWTAuth) GenerateToken(subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(a.duration)},
	})

	return token.SignedString([]byte(a.secret))
}

func (a *JWTAuth) AuthenticateRequest(r *http.Request) (string, error) {
	tokenStr := r.Header.Get("AUTH_TOKEN")
	if tokenStr == "" {
		return "", fmt.Errorf("no token provided")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf("invalid token")
}

func (a JWTAuth) Type() string {
	return "JWT"
}
