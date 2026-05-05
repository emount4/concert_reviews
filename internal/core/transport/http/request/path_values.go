package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"invalid path value, path value:'%s' by key:'%s': %w",
			pathValue,
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	pathValueInt, err := strconv.Atoi(pathValue)

	if err != nil {
		return 0, fmt.Errorf(
			"path value must be int path value:'%s' by key:'%s': %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return pathValueInt, nil
}
