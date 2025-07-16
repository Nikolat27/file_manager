package utils

import (
	"encoding/json"
	"io"
)

func ParseJsonData(reqBody io.ReadCloser, maxBytes int64, input any) error {
	body, err := io.ReadAll(io.LimitReader(reqBody, maxBytes))
	if err != nil {
		return err
	}

	return json.Unmarshal(body, input)
}
