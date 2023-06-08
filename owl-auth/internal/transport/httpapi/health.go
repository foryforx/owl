package httpapi

import (
	"context"
	"net/http"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/transport/httpapi/response"
)

type HealthChecker interface {
	CheckHealth(ctx context.Context) ([]api.Health, error)
}

type HealthCheckerFunc func(ctx context.Context) ([]api.Health, error)

func (h HealthCheckerFunc) CheckHealth(ctx context.Context) ([]api.Health, error) { return h(ctx) }

func HealthCheckHandler(h HealthChecker) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hh, err := h.CheckHealth(r.Context())
		if err != nil {
			response.RespondError(r.Context(), w, err, "Health check failed", http.StatusInternalServerError)
			return
		}

		response.RespondJSON(r.Context(), w, api.HealthResponse{Healths: hh}, http.StatusOK)
	}
}
