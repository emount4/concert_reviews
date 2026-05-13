package core_http_request

import (
	"fmt"
	"math"
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

func GetPageCount(total int, limit *int) int {
	if total <= 0 {
		return 0
	}
	if limit == nil || *limit <= 0 {
		return 1
	}
	return int(math.Ceil(float64(total) / float64(*limit)))
}
