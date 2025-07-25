package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
)

type UploadedFile struct {
	File multipart.File
	Size int64
	Name string
}

func ReadFile(r *http.Request, maxMemory int64, allowedFormats []string) (*UploadedFile, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	// getting the file (from http)
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	uploadedFile := &UploadedFile{
		File: file,
		Size: handler.Size,
		Name: handler.Filename,
	}

	if err := uploadedFile.validateFileFormat(allowedFormats); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}

func (file *UploadedFile) UploadToDisk(uploadDir string) (string, error) {
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

func (file *UploadedFile) validateFileFormat(allowedTypes []string) error {
	// Read first 512 bytes
	header := make([]byte, 512)
	n, err := file.File.Read(header)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(header[:n])
	if !slices.Contains(allowedTypes, contentType) {
		return fmt.Errorf("invalid file type: %s. Must be in :%s", contentType, allowedTypes)
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

func DeleteFileFromDisk(address string) error {
	return os.Remove(address)
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
