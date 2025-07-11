package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {}

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
