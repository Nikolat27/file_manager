package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
)

const (
	MaxMemoryBytes = 100
)

var allowedTypes = []string{"image/jpeg", "image/png", "application/zip"}

func UploadFileToDisk(key, userId string, r *http.Request) (string, error) {
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

	if err := validateFileFormat(file); err != nil {
		return "", err
	}

	uploadDir := "uploads/user_files/" + userId + "/"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	randomStr := rand.Text()
	fileName := randomStr + handler.Filename

	dst, err := os.Create(uploadDir + fileName)
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

func validateFileFormat(file multipart.File) error {
	// Read first 512 bytes
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(header[:n])
	if !slices.Contains(allowedTypes, contentType) {
		return errors.New("invalid file type: " + contentType)
	}

	// Reset file pointer to the beginning for further reading
	if seeker, ok := file.(io.Seeker); ok {
		if _, err = seeker.Seek(0, io.SeekStart); err != nil {
			return err
		}
	} else {
		return errors.New("cannot seek file")
	}

	return nil
}

func DeleteFileFromDisk(address []byte) error {
	err := os.Remove(string(address))
	if err != nil {
		return err
	}

	return nil
}
