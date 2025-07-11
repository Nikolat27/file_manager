package handlers

import (
	"encoding/json"
	"file_manager/utils"
	"fmt"
	"io"
	"net/http"
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

	if _, err := handler.Models.User.CreateUserInstance(input.Username, string(salt), string(hashedPassword[:])); err != nil {
		http.Error(w, fmt.Errorf("ERROR creating user instance: %s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("user created successfully"))
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		RawPassword string `json:"password"`
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

	user, err := handler.Models.User.FetchUserByUsername(input.Username)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR creating token: %s", err).Error(), http.StatusBadRequest)
		return
	}

	isPasswordTrue := utils.ValidateHash([]byte(input.RawPassword), []byte(user.HashedPassword), []byte(user.Salt))
	if !isPasswordTrue {
		http.Error(w, fmt.Errorf("ERROR validating password: %s", err).Error(), http.StatusBadRequest)
		return
	}
	
	//token, err := handler.PasetoMaker.CreateToken(input.Username, 24*time.Hour)
	//if err != nil {
	//	http.Error(w, fmt.Errorf("ERROR creating token: %s", err).Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//_, err = handler.PasetoMaker.VerifyToken(token)
	//if err != nil {
	//	http.Error(w, fmt.Errorf("ERROR verifying token: %s", err).Error(), http.StatusBadRequest)
	//	return
	//}

	w.Write([]byte("token is valid"))
}
