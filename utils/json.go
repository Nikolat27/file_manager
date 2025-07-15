package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseJsonData(r *http.Request, maxBytes int64, input any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBytes))
	if err != nil {
		return err
	}

	return json.Unmarshal(body, input)
}	
