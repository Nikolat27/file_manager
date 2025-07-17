package utils

import (
	"crypto/rand"
	"crypto/sha256"
)

const (
	Size       = 32
	SaltLength = 16
)

func Hash256(plainText, salt []byte) [Size]byte {
	newPassword := append(plainText, salt...)
	return sha256.Sum256(newPassword)
}

func ValidateHash(plainText, hashedText, salt []byte) bool {
	newText := append(plainText, salt...)
	currentHash := sha256.Sum256(newText)
	
	return string(currentHash[:]) == string(hashedText)
}

func GenerateSalt() ([]byte, error) {
	bytes := make([]byte, SaltLength)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
