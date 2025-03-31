package tools

import (
	"BTM-backend/internal/domain"
	"encoding/json"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseJWT(t *testing.T) {
	// Test data
	testData := map[string]interface{}{
		"user_id": "123",
		"email":   "test@example.com",
	}
	secret := "test-secret-key"

	// Generate token
	token, err := GenerateJWT(testData, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse token
	parsedData, err := ParseJWT(token, secret)
	require.NoError(t, err)
	assert.NotNil(t, parsedData)

	// Verify parsed data matches original
	var parsedMap map[string]interface{}
	err = json.Unmarshal(parsedData, &parsedMap)
	require.NoError(t, err)

	for key, value := range testData {
		assert.Equal(t, value, parsedMap[key])
	}
}

func TestParseInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	secret := "test-secret-key"

	_, err := ParseJWT(invalidToken, secret)
	assert.Error(t, err)
}

func TestParseExpiredToken(t *testing.T) {
	// Test data
	testData := map[string]interface{}{
		"user_id": "123",
		"email":   "test@example.com",
	}
	secret := "test-secret-key"

	// Generate token with expired time
	claims := jwt.MapClaims{
		"random": time.Now().UnixMicro(),
		"exp":    time.Now().Add(-24 * time.Hour).Unix(),
	}

	inrec, _ := json.Marshal(testData)
	err := json.Unmarshal(inrec, &claims)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	// Parse expired token
	_, err = ParseJWT(tokenString, secret)
	assert.Error(t, err)
	assert.Equal(t, "Token is expired", err.Error())
}

func TestFetchTokenInfo(t *testing.T) {
	// Create test context
	c := &gin.Context{}

	// Test data
	userInfo := domain.UserJwt{
		Account: "123",
		Role:    1,
	}

	// Test successful fetch
	c.Set("userInfo", userInfo)
	result, err := FetchTokenInfo(c)
	require.NoError(t, err)
	assert.Equal(t, userInfo, result)

	// Test missing userInfo
	c = &gin.Context{}
	_, err = FetchTokenInfo(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userInfo not found")

	// Test invalid type
	c.Set("userInfo", "invalid type")
	_, err = FetchTokenInfo(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse userInfo")
}

func TestParseToken(t *testing.T) {
	// Test data
	testData := domain.UserJwt{
		Account: "123",
		Role:    1,
	}
	secret := "test-secret-key"

	// Generate token
	token, err := GenerateJWT(testData, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse token
	result, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, testData, result)

	// Parse use error secret
	result, err = ParseToken(token, "error-secret")
	require.Error(t, err)
	assert.Equal(t, domain.UserJwt{}, result)
	assert.Contains(t, err.Error(), "ParseJWT")
}
