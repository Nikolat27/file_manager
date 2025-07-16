package utils

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func GetPaginationParams(r *http.Request) (int64, int64, error) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}

	intPage, err := strconv.ParseInt(page, 0, 64)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err := strconv.ParseInt(limit, 0, 64)
	if err != nil {
		return 0, 0, err
	}

	return intPage, intLimit, nil
}

func ReadFileId(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	fileId := params.ByName("file_id")
	if fileId == "" {
		return "", errors.New("file id is missing")
	}

	return fileId, nil
}

func ReadShortUrlParams(r *http.Request) ([]byte, error) {
	params := httprouter.ParamsFromContext(r.Context())
	fileId := params.ByName("file_short_url")
	if fileId == "" {
		return nil, errors.New("file short url is missing")
	}

	return []byte(fileId), nil
}
