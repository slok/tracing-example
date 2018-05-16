package server

import (
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

func (s *server) middlewareLogger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.With("method", r.Method).With("URI", r.RequestURI).Infof("request received")
		next.ServeHTTP(w, r)
	})
}

func (s *server) middlewareTrace(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := &responseWriterInterceptor{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Get context from ingress request.
		spCtx, err := s.tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			s.logger.Warnf("no tracing on headers")
			next.ServeHTTP(iw, r)
		}

		// Create a new span.
		span := s.tracer.StartSpan("api-request", opentracing.ChildOf(spCtx))
		defer span.Finish()
		// Set the final result after executing all the request chain.
		defer func(start time.Time) {
			span.SetTag("span.kind", "server")
			span.SetTag("http.method", r.Method)
			span.SetTag("http.status_code", iw.status)
			span.SetTag("http.url", r.URL)
			span.SetTag("app.layer", "middleware")

			span.LogKV(
				"remote_addr", r.RemoteAddr,
				"method", r.Method,
				"url", r.URL,
				"content_length", r.ContentLength,
				"status_code", iw.status,
				"status_text", http.StatusText(iw.status),
				"response_size", iw.size,
				"took", time.Since(start).String(),
				"sec", time.Since(start).Seconds(),
			)
		}(time.Now())

		// update request context for the new ones.
		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)
		next.ServeHTTP(iw, r)
	})
}
