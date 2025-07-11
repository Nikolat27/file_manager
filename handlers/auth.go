package handlers

import (
	"encoding/json"
	"file_manager/utils"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := utils.ReadJson(r, 10000, &input); err != nil {
		http.Error(w, fmt.Sprintf("ERROR reading json: %s", err), http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "ERROR both username and password are required", http.StatusBadRequest)
		return
	}

	salt, err := utils.GetSalt()
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR create salt: %s", err).Error(), http.StatusBadRequest)
		return
	}

	hashedPassword := utils.Hash256([]byte(input.Password), salt)

	fmt.Println(hashedPassword)
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 10000))
	if err != nil {
		fmt.Println("Error Reader: ", err)
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		fmt.Println("Error UnMarshaller: ", err)
		return
	}

	token, err := handler.PasetoMaker.CreateToken(input.Username, 24*time.Hour)
	if err != nil {
		fmt.Println(err)
		return
	}

	payload, err := handler.PasetoMaker.VerifyToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(payload.Valid())
}
