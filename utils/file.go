package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	MaxMemoryBytes = 100
)

func UploadFile(key, userId string, r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(MaxMemoryBytes << 20); err != nil {
		return "", err
	}

	file, handler, err := r.FormFile(key)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	uploadDir := "uploads/user_files/" + userId + "/"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(uploadDir + handler.Filename)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := dst.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()
	
	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return dst.Name(), nil
}
