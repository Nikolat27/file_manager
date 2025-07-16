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

	if err := utils.ParseJSON(r.Body, 10000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading json: %w", err))
		return
	}

	if input.Username == "" || input.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("both username and password are required"))
		return
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating salt: %w", err))
		return
	}

	hashedPassword := utils.Hash256([]byte(input.Password), salt)
	encodedHash := hex.EncodeToString(hashedPassword[:])
	encodedSalt := hex.EncodeToString(salt)

	_, err = handler.Models.User.GetByUsername(input.Username)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("this username is taken already"))
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("fetch user: %w", err))
		return
	}

	_, err = handler.Models.User.Create(input.Username, encodedSalt, encodedHash)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("creating user instance: %w", err))
		return
	}

	utils.WriteJSON(w, "user registered successfully")
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		RawPassword string `json:"password"`
	}

	if err := utils.ParseJSON(r.Body, 10000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading json: %w", err))
		return
	}

	user, err := handler.Models.User.GetByUsername(input.Username)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		return
	}

	decodedHash, err := hex.DecodeString(user.HashedPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error decoding hash: %w", err))
		return
	}

	decodedSalt, err := hex.DecodeString(user.Salt)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error decoding salt: %w", err))
		return
	}

	if !utils.ValidateHash([]byte(input.RawPassword), decodedHash, decodedSalt) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid username or password"))
		return
	}

	token, err := handler.PasetoMaker.CreateToken(user.Username, user.Id.Hex(), 24*time.Hour)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating token: %w", err))
		return
	}

	utils.WriteJSON(w, token)
}
