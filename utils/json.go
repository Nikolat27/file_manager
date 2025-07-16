package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseJSON(reqBody io.ReadCloser, maxBytes int64, input any) error {
	body, err := io.ReadAll(io.LimitReader(reqBody, maxBytes))
	if err != nil {
		return err
	}

	return json.Unmarshal(body, input)
}

func WriteJSON[T interface{ []byte | string }](w http.ResponseWriter, msg T) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(msg)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func WriteError(w http.ResponseWriter, status int, err any) {
	switch val := err.(type) {
	case string:
		http.Error(w, val, status)
	case error:
		http.Error(w, val.Error(), status)
	}
}
