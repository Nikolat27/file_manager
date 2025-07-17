package token

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

const symmetricKey = "$5rjt5n^euch!hwqar2%p$uydrgtnh%a"

func New() (*PasetoMaker, error) {
	key, err := getSymmetricKey()
	if err != nil {
		return nil, err
	}

	if len(key[:]) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: key[:],
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username, userId, userPlan string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, userId, userPlan, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil); err != nil {
		return nil, err
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}

func getSymmetricKey() ([32]byte, error) {
	if symmetricKey == "" {
		return [32]byte{}, errors.New("symmetric key is empty")
	}

	hash := sha256.Sum256([]byte(symmetricKey))

	return hash, nil
}
