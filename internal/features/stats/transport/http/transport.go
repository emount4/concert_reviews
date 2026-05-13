package stats_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

type StatsHTTPHandler struct {
	statsService StatsService
}

type StatsService interface {
	GetGlobalStats(ctx context.Context) (domain.GlobalStats, error)
}

func NewStatsHTTPHandler(statsService StatsService) *StatsHTTPHandler {
	return &StatsHTTPHandler{
		statsService: statsService,
	}
}

func (h *StatsHTTPHandler) GetGlobalStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	stats, err := h.statsService.GetGlobalStats(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get stats")
		return
	}

	responseHandler.JSONResponse(MapDomainToGlobalStatsResponse(stats), http.StatusOK)
}

func (h *StatsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/stats/global",
			Handler: h.GetGlobalStats,
		},
	}
}
