package middleware

import (
	"net/http"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
)

type LogRequest struct {
	log logger.Logger
}

func NewLogRequest(log logger.Logger) *LogRequest {
	return &LogRequest{
		log: log,
	}
}

func (m LogRequest) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dw := newResponseDataWriter(w)

		start := time.Now()
		next.ServeHTTP(dw, r)
		duration := time.Since(start)

		m.log.Info("Got request", map[string]interface{}{
			"uri":      r.URL.Path,
			"method":   r.Method,
			"duration": duration,
			"status":   dw.data.status,
			"size":     dw.data.size,
		})
	})
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
