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

func WriteJSONData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var message string
	switch val := err.(type) {
	case string:
		message = val
	case error:
		message = val.Error()
	default:
		message = "unknown error"
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"error": message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
