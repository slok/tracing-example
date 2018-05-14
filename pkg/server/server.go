package server

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/slok/tracing-example/pkg/log"
)

type server struct {
	router *http.ServeMux
	// endpoints are the services that will be called randomly
	endpoints []string
	logger    log.Logger
	tracer    opentracing.Tracer
}

// New returns a new server.
func New(endpoints []string, tracer opentracing.Tracer, mux *http.ServeMux, logger log.Logger) http.Handler {
	s := &server{
		endpoints: endpoints,
		tracer:    tracer,
		router:    mux,
		logger:    logger,
	}

	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
