package utils

import (
	"errors"
	"file_manager/token"
	"net/http"
)

func CheckAuth(r *http.Request, paseto *token.PasetoMaker) (*token.Payload, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return nil, errors.New("unauthorized: authToken is missing")
	}

	payload, err := paseto.VerifyToken(authToken)
	if err != nil {
		return nil, errors.New("unauthorized: invalid token")
	}

	return payload, nil
}
