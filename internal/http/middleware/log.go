package middleware

import (
	"net/http"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
)

// LogMiddleware HTTP middleware логирует запрос.
func LogMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dw := newResponseDataWriter(w)

			start := time.Now()
			next.ServeHTTP(dw, r)
			duration := time.Since(start)

			logger.
				WithField("uri", r.URL.Path).
				WithField("method", r.Method).
				WithField("content_type", r.Header.Get("Content-Type")).
				WithField("duration", duration).
				WithField("status", dw.data.status).
				WithField("size", dw.data.size).
				Info("Got request")
		}

		return http.HandlerFunc(fn)
	}
}

type responseData struct {
	status int
	size   int
}

type responseDataWriter struct {
	http.ResponseWriter

	data *responseData
}

func newResponseDataWriter(w http.ResponseWriter) *responseDataWriter {
	return &responseDataWriter{
		ResponseWriter: w,
		data: &responseData{
			status: http.StatusOK,
			size:   0,
		},
	}
}

func (w *responseDataWriter) Write(buf []byte) (int, error) {
	size, err := w.ResponseWriter.Write(buf)
	w.data.size += size
	return size, err
}

func (w *responseDataWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.data.status = status
}
