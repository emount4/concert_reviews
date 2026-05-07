package core_http_request

import (
	"fmt"
	"net/http"
)

func GetLimitOffsetByQueryParam(r *http.Request) (*int, *int, error) {
	limit, err := GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
