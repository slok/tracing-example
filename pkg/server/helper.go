package server

import (
	"net/http"
)

// responseWriterInterceptor is an incerceptor so the middleware can inspect the
// setted status codes.
type responseWriterInterceptor struct {
	size   int
	status int
	http.ResponseWriter
}

func (w *responseWriterInterceptor) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriterInterceptor) Write(p []byte) (int, error) {
	w.size += len(p)
	return w.ResponseWriter.Write(p)
}
