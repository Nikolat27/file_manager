package webserver

import (
	"errors"
	"file_manager/handlers"
	"net/http"
)

func CheckAuth(handler *handlers.Handler, httpHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := getToken(w, r, "authToken")
		if err != nil {
			return
		}

		if _, err := handler.PasetoMaker.VerifyToken(value); err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		httpHandler(w, r)
	}
}

func getToken(w http.ResponseWriter, r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "UnAuthorized: Cookie not found", http.StatusUnauthorized)
			return "", err
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return "", err
	}

	return cookie.Value, nil
}
