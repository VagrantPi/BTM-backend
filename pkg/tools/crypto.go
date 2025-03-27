package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const PwdCost = 12

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PwdCost)
	return string(bytes), err
}

// CheckPassword 這裡慢是正常的，https://stackoverflow.com/a/49437553
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CheckPasswordRule 檢查密碼是否符合規定
func CheckPasswordRule(password string) bool {
	return len(password) >= 10 &&
		// 至少包含一個小寫字母
		regexp.MustCompile(`[a-z]`).MatchString(password) &&
		// 至少包含一個大寫字母
		regexp.MustCompile(`[A-Z]`).MatchString(password) &&
		// 至少包含一個數字
		regexp.MustCompile(`[0-9]`).MatchString(password) &&
		// 至少包含一個特殊字符
		regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
}

// 機敏資料單向Hash
func HashSensitiveData(key, data string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("server key is empty")
	}
	hash := hmac.New(sha256.New, []byte(key))
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// EncryptAES256 encrypts a string using key
func EncryptAES256(key, data string) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: %d", len(key))
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(data), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// DecryptAES256 decrypts a string using key
func DecryptAES256(key, data string) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: %d", len(key))
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}

// MaskName 會將字串進行遮罩處理
func MaskName(name string) string {
	// 使用 utf8.RuneCountInString 來正確計算 Unicode 字元數
	nameLen := utf8.RuneCountInString(name)

	// 將名字轉換為 rune 切片，以正確處理 Unicode 字元
	nameRunes := []rune(name)

	if nameLen == 2 {
		// 對於兩個字的名字，保留第一個字，第二個字遮罩
		return string(nameRunes[0:1]) + "X"
	}

	if nameLen > 2 {
		// 對於超過兩個字的名字，保留第一個和最後一個字，中間的字遮罩
		maskedMiddle := make([]rune, nameLen-2)
		for i := range maskedMiddle {
			maskedMiddle[i] = 'X'
		}
		return string(nameRunes[0:1]) + string(maskedMiddle) + string(nameRunes[nameLen-1:])
	}

	return name
}

// MaskEmail 會將email進行遮罩處理
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	atIndex := strings.Index(email, "@")
	if atIndex <= 5 {
		return email[:atIndex-1] + "x" + email[atIndex:]
	}

	return email[:5] + strings.Repeat("x", atIndex-5) + email[atIndex:]
}
