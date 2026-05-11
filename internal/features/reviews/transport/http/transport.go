package review_transport_http

import core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"

type ReviewService interface {
}

type ReviewHTTPHandler struct {
	reviewService ReviewService
}

func NewReviewHTTPHandler(s ReviewService) *ReviewHTTPHandler {
	return &ReviewHTTPHandler{
		reviewService: s,
	}
}

func (h *ReviewHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{}
}
