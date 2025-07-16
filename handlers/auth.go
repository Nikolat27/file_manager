package handlers

import (
	"encoding/hex"
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := utils.ParseJson(r.Body, 10000, &input); err != nil {
		http.Error(w, fmt.Sprintf("ERROR reading json: %s", err), http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "ERROR both username and password are required", http.StatusBadRequest)
		return
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR create salt: %s", err).Error(), http.StatusBadRequest)
		return
	}

	hashedPassword := utils.Hash256([]byte(input.Password), salt)

	encodedHash := hex.EncodeToString(hashedPassword[:])
	encodedSalt := hex.EncodeToString(salt)

	if _, err = handler.Models.User.GetByUsername(input.Username); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := handler.Models.User.Create(input.Username, encodedSalt, encodedHash); err != nil {
				http.Error(w, fmt.Errorf("creating user instance: %s", err).Error(), http.StatusBadRequest)
				return
			}

			return
		}
		http.Error(w, fmt.Errorf("fetch user: %s", err).Error(), http.StatusBadRequest)
		return
	}

	http.Error(w, "this username is taken already", http.StatusBadRequest)
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		RawPassword string `json:"password"`
	}

	if err := utils.ParseJson(r.Body, 10000, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := handler.Models.User.GetByUsername(input.Username)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR creating token: %s", err).Error(), http.StatusBadRequest)
		return
	}

	decodedHash, err := hex.DecodeString(user.HashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decodedSalt, err := hex.DecodeString(user.Salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.ValidateHash([]byte(input.RawPassword), decodedHash, decodedSalt) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	token, err := handler.PasetoMaker.CreateToken(user.Username, user.Id.Hex(), 24*time.Hour)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR creating token: %s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(token))
}
