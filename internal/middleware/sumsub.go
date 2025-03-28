package middleware

import (
	"BTM-backend/configs"
	"BTM-backend/pkg/error_code"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

var SumsubGuardImpl SumsubGuard

type SumsubGuard struct{}

func (m SumsubGuard) CheckDigest(c *gin.Context) {
	sumsubPrivateKey := configs.C.Sumsub.WebhookSecretKey
	algoHeader := c.GetHeader("X-Payload-Digest-Alg")
	hashFunc, err := m.getAlgMethod(algoHeader)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	bodyBytes, err := m.copyBody(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	mac := hmac.New(hashFunc, []byte(sumsubPrivateKey))
	mac.Write(bodyBytes)
	expectedMAC := mac.Sum(nil)
	providedDigest := c.GetHeader("X-Payload-Digest")
	calculatedDigest := hex.EncodeToString(expectedMAC)

	if err = m.validateDigest(providedDigest, calculatedDigest); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}

func (m SumsubGuard) getAlgMethod(algoHeader string) (hashFunc func() hash.Hash, err error) {
	switch algoHeader {
	case "HMAC_SHA1_HEX":
		hashFunc = sha1.New
	case "HMAC_SHA256_HEX":
		hashFunc = sha256.New
	case "HMAC_SHA512_HEX":
		hashFunc = sha512.New
	default:
		msg := fmt.Sprintf("unsupported algorithm: %s", algoHeader)
		err = errors.BadRequest(error_code.ErrSumsubBadRequest, msg)
		return nil, err
	}
	return hashFunc, nil
}

func (m SumsubGuard) copyBody(c *gin.Context) (bodyBytes []byte, err error) {
	if c.Request.Body != nil {
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			err = errors.BadRequest(error_code.ErrSumsubApiUnmarshal, "io.ReadAll(reader)").WithCause(err)
			return bodyBytes, err
		}
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return
}

func (m SumsubGuard) validateDigest(providedDigest string, calculatedDigest string) (err error) {
	if !hmac.Equal([]byte(calculatedDigest), []byte(providedDigest)) {
		err = errors.BadRequest(error_code.ErrSumsubApiValidate, "Bad Digest").WithMetadata(map[string]string{
			"providedDigest":   providedDigest,
			"calculatedDigest": calculatedDigest,
		})
		return err
	}
	return
}
