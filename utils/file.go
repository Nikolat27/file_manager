package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
)

var allowedTypes = []string{"image/jpeg", "image/png", "application/zip"}

type UploadedFile struct {
	File multipart.File
	Size int64
	Name string
}

func ReadFile(r *http.Request, maxMemory int64) (*UploadedFile, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	uploadedFile := &UploadedFile{
		File: file,
		Size: handler.Size,
		Name: handler.Filename,
	}

	if err := uploadedFile.validateFileFormat(); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}

func (file *UploadedFile) UploadToDisk(userId string) (string, error) {
	uploadDir := "uploads/user_files/" + userId + "/"
	if err := makeDir(uploadDir); err != nil {
		return "", err
	}

	dst, err := os.Create(uploadDir + file.generateFileName())
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file.File); err != nil {
		return "", err
	}

	fileAddress := dst.Name()
	return fileAddress, nil
}

func (file *UploadedFile) generateFileName() string {
	randomStr := rand.Text()
	return randomStr + file.Name
}

func (file *UploadedFile) validateFileFormat() error {
	// Read first 512 bytes
	header := make([]byte, 512)
	n, err := file.File.Read(header)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(header[:n])
	if !slices.Contains(allowedTypes, contentType) {
		return errors.New("invalid file type: " + contentType)
	}

	// Reset file pointer to the beginning for further reading
	if seeker, ok := file.File.(io.Seeker); ok {
		if _, err = seeker.Seek(0, io.SeekStart); err != nil {
			return err
		}
	} else {
		return errors.New("cannot seek file")
	}

	return nil
}

func makeDir(uploadDir string) error {
	return os.MkdirAll(uploadDir, 0755)
}

func DeleteFileFromDisk(address []byte) error {
	return os.Remove(string(address))
}

func CheckFilePassword(hashedPassword, salt, rawPassword []byte) error {
	decodedHashPassword, err := hex.DecodeString(string(hashedPassword))
	if err != nil {
		return err
	}

	decodedSalt, err := hex.DecodeString(string(salt))
	if err != nil {
		return err
	}

	if !ValidateHash(rawPassword, decodedHashPassword, decodedSalt) {
		return errors.New("password is invalid")
	}

	return nil
}
