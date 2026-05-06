package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)

	if err != nil {
		return nil, fmt.Errorf(
			"params='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return &val, nil
}

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	return &param
}

func GetBoolQueryParam(r *http.Request, key string) *bool {
	param := r.URL.Query().Get(key)

	switch param {
	case "true":
		v := true
		return &v
	case "false":
		v := false
		return &v
	default:
		return nil
	}
}
