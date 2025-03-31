package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveFile(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Remove existing file
	err = RemoveFile(tempFile.Name())
	require.NoError(t, err)
	assert.False(t, FileExists(tempFile.Name()))

	// Attempt to remove non-existent file
	err = RemoveFile(tempFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCreateFile(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "testfile")
	defer os.Remove(tempFile)

	// Test creating file with data
	data := []byte("test data")
	err := CreateFile(tempFile, data)
	require.NoError(t, err)

	// Verify file exists and contains correct data
	assert.True(t, FileExists(tempFile))
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	assert.Equal(t, data, content)

	// Test creating file with empty data
	err = CreateFile(tempFile, nil)
	require.NoError(t, err)
	content, err = os.ReadFile(tempFile)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestFileExists(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Test existing file
	assert.True(t, FileExists(tempFile.Name()))

	// Test non-existent file
	assert.False(t, FileExists("nonexistent_file.txt"))

	// Test non-existent file with relative path
	assert.False(t, FileExists("../nonexistent_file.txt"))

	// Test directory
	tempDir := filepath.Join(os.TempDir(), "testdir")
	err = os.Mkdir(tempDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	assert.True(t, FileExists(tempDir))
}
