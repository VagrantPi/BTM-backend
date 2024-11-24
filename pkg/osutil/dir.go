package osutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// FilePath 使用方式：根目錄的 ("dir", "dir2", "fileName")
// 由於 golang file path 是根據執行的資料夾決定的，因此在 go run, go build 也會有不同場景
func FilePath(pathAndFilename ...string) (p string, err error) {
	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	WorkPath, err := os.Getwd()
	if err != nil {
		return
	}

	p = filepath.Join(WorkPath, filepath.Join(pathAndFilename...))
	if !FileExists(p) {
		p = filepath.Join(appPath, filepath.Join(pathAndFilename...))
		if !FileExists(p) {
			return "", fmt.Errorf("file not found: %s", p)
		}
	}

	return
}
