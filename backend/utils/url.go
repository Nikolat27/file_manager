package utils

import (
	"context"
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

func ParseIdParam(requestContext context.Context) (string, error) {
	params := httprouter.ParamsFromContext(requestContext)
	id := params.ByName("id")
	if id == "" {
		return "", errors.New("id is missing")
	}

	return id, nil
}
