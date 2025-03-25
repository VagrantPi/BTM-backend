package tools

import (
	"regexp"

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
