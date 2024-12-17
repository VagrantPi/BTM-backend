package tools

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/error_code"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT Deprecated: 注意不要放敏感資訊
func GenerateJWT(v any, secret string) (string, error) {
	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	hmacSampleSecret := []byte(secret)

	var claims = jwt.MapClaims{
		"random": time.Now().UnixMicro(),
		"exp":    time.Now().Add(3 * 24 * time.Hour).Unix(),
	}
	inrec, _ := json.Marshal(v)
	err := json.Unmarshal(inrec, &claims)
	if err != nil {
		return "", err
	}

	// For HMAC signing method, the key can be any []byte. It is recommended to generate
	// a key using crypto/rand or something equivalent. You need the same key for signing
	// and validating.

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(hmacSampleSecret)
}

// ParseJWT
func ParseJWT(tokenString string, secret string) ([]byte, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		inrec, err := json.Marshal(claims)
		if err != nil {
			return nil, err
		}
		return inrec, nil
	}

	return nil, errors.Unauthorized(error_code.ErrInvalidJWT, "ParseJWT")
}

func ParseToken(token string) (claim domain.UserJwt, err error) {
	data, err := ParseJWT(token, configs.C.JWT.Secret)
	if err != nil {
		err = errors.Unauthorized(error_code.ErrInvalidJWTParse, "ParseJWT").WithCause(err)
		return
	}

	err = json.Unmarshal(data, &claim)
	if err != nil {
		err = errors.Unauthorized(error_code.ErrInvalidJWT, "Unmarshal").WithCause(err)
		return
	}

	return
}
