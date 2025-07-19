package webserver

import (
	"errors"
	"file_manager/handlers"
	"net/http"
	"os"
	"strings"
)

// not going to use it for a while...

//func CheckAuth(handler *handlers.Handler, httpHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		value, err := getToken(w, r, "authToken")
//		if err != nil {
//			return
//		}
//
//		if _, err := handler.PasetoMaker.VerifyToken(value); err != nil {
//			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		httpHandler(w, r)
//	}
//}

func CheckAuth(handler *handlers.Handler, httpHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized: token is missing", http.StatusUnauthorized)
			return
		}

		if _, err := handler.PasetoMaker.VerifyToken(token); err != nil {
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

func CORSMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getAllowedCORS()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if allowed == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				break
			}

			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", allowed)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// browser preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAllowedCORS() []string {
	origins := os.Getenv("ALLOWED_CORS_ORIGINS")
	if origins == "" {
		return []string{"*"} // default origin
	}

	return strings.Split(origins, ",")
}
