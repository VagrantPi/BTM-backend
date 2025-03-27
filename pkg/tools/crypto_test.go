package tools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "正常情況",
			password: "test123",
			wantErr:  false,
		},
		{
			name:     "空密碼",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GeneratePasswordHash(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "test123"
	hash, err := GeneratePasswordHash(password)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "正確密碼",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "錯誤密碼",
			password: "wrong",
			hash:     hash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.hash, tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckPasswordRule(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "符合規定",
			password: "Test123!@#",
			want:     true,
		},
		{
			name:     "太短",
			password: "Test123",
			want:     false,
		},
		{
			name:     "缺少小寫字母",
			password: "TEST123!@#",
			want:     false,
		},
		{
			name:     "缺少大寫字母",
			password: "test123!@#",
			want:     false,
		},
		{
			name:     "缺少數字",
			password: "TestABC!@#",
			want:     false,
		},
		{
			name:     "缺少特殊字符",
			password: "Test123ABC",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordRule(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashSensitiveData(t *testing.T) {
	key := "testkey12345678901234567890123456789012" // 32字元的key
	tests := []struct {
		name    string
		data    string
		key     string
		wantErr bool
	}{
		{
			name:    "正常情況",
			data:    "test data",
			key:     key,
			wantErr: false,
		},
		{
			name:    "空key",
			data:    "test data",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashSensitiveData(tt.key, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func TestEncryptDecryptAES256(t *testing.T) {
	key := "testkey1234567890123456789012345" // 32字元的key
	tests := []struct {
		name    string
		data    string
		key     string
		wantErr bool
	}{
		{
			name:    "正常情況",
			data:    "test data",
			key:     key,
			wantErr: false,
		},
		{
			name:    "空key",
			data:    "test data",
			key:     "",
			wantErr: true,
		},
		{
			name:    "key長度過短，但因為有補齊 key 所以應正確",
			data:    "test data",
			key:     "short",
			wantErr: false,
		},
		{
			name:    "key長度過長，但因為有裁剪 key 所以應正確",
			data:    "test data",
			key:     key + "12345678901234567890123456789012",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptAES256(tt.key, tt.data)
			fmt.Println("err", err)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			decrypted, err := DecryptAES256(tt.key, encrypted)
			assert.NoError(t, err)
			assert.Equal(t, tt.data, decrypted)
		})
	}
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字串",
			input:    "",
			expected: "",
		},
		{
			name:     "兩個字",
			input:    "張三",
			expected: "張X",
		},
		{
			name:     "三個字",
			input:    "張三豐",
			expected: "張X豐",
		},
		{
			name:     "四個字",
			input:    "張三豐四",
			expected: "張XX四",
		},
		{
			name:     "一個字",
			input:    "張",
			expected: "張",
		},
		{
			name:     "英文名字",
			input:    "John",
			expected: "JXXn",
		},
		{
			name:     "特殊字符",
			input:    "@#%&",
			expected: "@XX&",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskName(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "空地址",
			email:    "",
			expected: "",
		},
		{
			name:     "短地址，@前小於5個字元",
			email:    "a@b.com",
			expected: "x@b.com",
		},
		{
			name:     "短地址，@前等於5個字元",
			email:    "abcde@b.com",
			expected: "abcdx@b.com",
		},
		{
			name:     "長地址，@前大於5個字元",
			email:    "abcdefg@b.com",
			expected: "abcdexx@b.com",
		},
		{
			name:     "正常地址",
			email:    "test.email@example.com",
			expected: "test.xxxxx@example.com",
		},
		{
			name:     "特殊字符",
			email:    "user+test@domain.com",
			expected: "user+xxxx@domain.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskEmail(tt.email)
			if got != tt.expected {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.email, got, tt.expected)
			}
		})
	}
}
