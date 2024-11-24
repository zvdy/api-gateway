package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zvdy/api-gateway/internal/auth"
)

func TestGenerateJWT(t *testing.T) {
	token, err := auth.GenerateJWT("test-user-id")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	// Generate a token
	token, err := auth.GenerateJWT("test-user-id")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := auth.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check the claims
	userID, ok := claims["user_id"].(string)
	assert.True(t, ok)
	assert.Equal(t, "test-user-id", userID)

	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.True(t, time.Unix(int64(exp), 0).After(time.Now()))
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	// Validate an invalid token
	_, err := auth.ValidateJWT("invalid-token")
	assert.Error(t, err)
}
