package user_transport_http

import (
	"encoding/json"
	"net/http"
)

type MeResponse struct {
	Message string `json:"message"`
}

func (h *UserHTTPHandler) Me(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	response := MeResponse{Message: "user profile endpoint is not implemented yet"}
	_ = json.NewEncoder(rw).Encode(response)
}
