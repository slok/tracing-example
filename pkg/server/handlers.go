package server

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

const (
	fastEndMin = 10
	fastEndMax = 100

	slowEndMin = 600
	slowEndMax = 3000

	multipleCallNumberMin = 1
	multipleCallNumberMax = 4
)

func (s *server) fastEnd() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dur := s.fastRandomDuration()
		time.Sleep(dur)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("service fast end"))
	})
}

func (s *server) slowEnd() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dur := s.slowRandomDuration()
		time.Sleep(dur)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("service slow end"))
	})
}

func (s *server) singleCall() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.endpoints) == 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("no endpoints available"))
			return
		}

		endpoint := s.getRandomEndpointWithRandomPath()
		req, err := s.newTracedRequest(r, "GET", endpoint, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error making request to %s: %s", endpoint, err)
		}
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error making request to %s: %s", endpoint, err)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("service single call"))
	})
}

func (s *server) multipleCalls() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.endpoints) == 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("no endpoints available"))
			return
		}

		// Make random multiple calls.
		times := s.randomNumber(multipleCallNumberMin, multipleCallNumberMax)
		for i := 0; i < times; i++ {
			endpoint := s.getRandomEndpointWithRandomPath()
			req, err := s.newTracedRequest(r, "GET", endpoint, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error making request to %s: %s", endpoint, err)
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error making request to %s: %s", endpoint, err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("service multiple call"))
	})
}

func (s *server) fastRandomDuration() time.Duration {

	return s.randomDuration(fastEndMin, fastEndMax)
}

func (s *server) slowRandomDuration() time.Duration {
	return s.randomDuration(slowEndMin, slowEndMax)
}

func (s *server) randomDuration(minMS, maxMS int) time.Duration {
	randNum := s.randomNumber(minMS, maxMS)
	return time.Duration(randNum) * time.Millisecond
}

func (s *server) getRandomEndpointWithRandomPath() string {
	src := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(src)
	endIdx := r.Intn(len(s.endpoints))
	routeIdx := r.Intn(len(routePaths))

	return fmt.Sprintf("%s%s", s.endpoints[endIdx], routePaths[routeIdx])
}

func (s *server) randomNumber(minMS, maxMS int) int {
	src := rand.NewSource(int64(time.Now().Nanosecond()))
	return rand.New(src).Intn(maxMS-minMS) + minMS
}

func (s *server) newTracedRequest(srcReq *http.Request, method, url string, body io.Reader) (*http.Request, error) {
	var spanCtx opentracing.SpanContext

	// Get current span.
	span := opentracing.SpanFromContext(srcReq.Context())
	// If no span present on context then get from the headers
	if span == nil {
		// Get context from source request.
		var err error
		spanCtx, err = s.tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(srcReq.Header))
		if err != nil {
			return nil, fmt.Errorf("could not extract opentracing headers: %s", err)
		}
	} else {
		spanCtx = span.Context()
	}
	// Create new request.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Inject context on dst request.
	s.tracer.Inject(
		spanCtx,
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	return req, nil
}
