package tools

import (
	"golang.org/x/crypto/bcrypt"
)

const PwdCost = 12

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PwdCost)
	return string(bytes), err
}

// CheckPassword 這裡慢是正常的，https://stackoverflow.com/a/49437553
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}
