package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"eth-fetcher/auth"
)

func TestJWTAuth_AuthenticateRequest_ValidToken(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, time.Hour)

	// Create a valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	// Create a request with the token in the header
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", tokenStr)

	// Call the AuthenticateRequest function
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that the authentication was successful
	assert.NoError(t, err)
	assert.Equal(t, "user123", sub)
}

func TestJWTAuth_AuthenticateRequest_NoTokenProvided(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, time.Hour)

	// Create a request without a token in the header
	req, _ := http.NewRequest("GET", "/", nil)

	// Call the AuthenticateRequest function
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.EqualError(t, err, "no token provided")
	assert.Empty(t, sub)
}

func TestJWTAuth_AuthenticateRequest_InvalidToken(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, time.Hour)

	// Create an invalid token
	invalidTokenStr := "invalid-token"

	// Create a request with the invalid token in the header
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", invalidTokenStr)

	// Call the AuthenticateRequest function
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Empty(t, sub)
}

func TestJWTAuth_AuthenticateRequest_ExpiredToken(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, 0)

	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": 1234567890, // Set the expiration time to a past timestamp
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	// Create a request with the expired token in the header
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", tokenStr)

	// Call the AuthenticateRequest function
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Empty(t, sub)
}

func TestJWTAuth_AuthenticateRequest_MalformedToken(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, 0)

	// Create a malformed token
	malformedTokenStr := "malformed-token"

	// Create a request with the malformed token in the header
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", malformedTokenStr)

	// Call the AuthenticateRequest function
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Empty(t, sub)
}

func TestJWTAuth_AuthenticateRequest_ErrorHandling(t *testing.T) {
	// Create a JWTAuth instance with a secret
	secret := "my-secret"
	authenticator := auth.NewJWTAuth(secret, 0)

	// Create a request with an unexpected error
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", "valid-token")

	// Call the AuthenticateRequest function with an error
	sub, err := authenticator.AuthenticateRequest(req)

	// Assert that the error is handled correctly
	assert.Error(t, err)
	assert.Empty(t, sub)
}
