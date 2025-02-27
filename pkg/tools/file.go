package tools

import (
	"BTM-backend/pkg/error_code"
	"io"
	"log"
	"os"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/yeka/zip"
)

func RemoveFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return errors.InternalServer(error_code.ErrRemoveFile, "remove file err").WithCause(err)
	}
	return nil
}

func CreateFile(filename string, data []byte) error {
	// 不管存不存在，先刪除，以避免檔案已存在
	err := RemoveFile(filename)
	if err != nil && os.IsNotExist(err) {
		return errors.InternalServer(error_code.ErrRemoveFile, "remove old file err").WithCause(err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return errors.InternalServer(error_code.ErrCreateFile, "create file err").WithCause(err)
	}

	return nil
}

func UnzipFile(zipPath, destFile, pwd string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.InternalServer(error_code.ErrToolsUnzipFile, "open zip file err").WithCause(err)
	}
	defer r.Close()

	if len(r.File) == 0 {
		return errors.InternalServer(error_code.ErrToolsUnzipFile, "zip file is empty")
	}

	f := r.File[0]
	if f.IsEncrypted() && pwd != "" {
		f.SetPassword(pwd)
	}

	r2, err := f.Open()
	if err != nil {
		log.Fatal(err)
	}

	buf, err := io.ReadAll(r2)
	if err != nil {
		log.Fatal(err)
	}
	defer r2.Close()

	err = CreateFile(destFile, buf)
	if err != nil {
		return errors.InternalServer(error_code.ErrToolsUnzipFile, "create file err").WithCause(err)
	}

	return nil
}
