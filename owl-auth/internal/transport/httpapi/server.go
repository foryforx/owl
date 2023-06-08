package httpapi

import (
	"net/http"

	"github.com/foryforx/owl/owl-auth/internal/conf"
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ConnPool      *sqlx.DB
	Conf          *conf.Configuration
	HealthChecker HealthChecker
	Middleware    Middleware
	Authenticator Middleware
}

func (s *APIServer) CreateMux(routers ...Router) *mux.Router {
	serveMux := mux.NewRouter()

	for _, router := range routers {
		for _, route := range router.Routes() {
			h := s.Middleware(route.Handler)
			if !route.IsPubic {
				h = s.Authenticator(h)
			}
			serveMux.Methods(route.Method).Path(route.Path).Handler(ToStdHandler(h))
		}
	}
	serveMux.Methods(http.MethodGet).Path("/health").Handler(ToStdHandler(HealthCheckHandler(s.HealthChecker)))

	return serveMux
}
