package tools

import (
	"fmt"
	"runtime"
)

// OPTIMIZE: 動態的判斷caller層數
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
