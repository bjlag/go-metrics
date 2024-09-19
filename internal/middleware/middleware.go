package middleware

import (
	"net/http"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
)

type responseData struct {
	status int
	size   int
}

type responseDataWriter struct {
	http.ResponseWriter

	data *responseData
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

func CreateLogRequestMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &responseDataWriter{
				ResponseWriter: w,
				data: &responseData{
					status: 0,
					size:   0,
				},
			}

			start := time.Now()
			next.ServeHTTP(lw, r)
			duration := time.Since(start)

			logger.Info("Got request", map[string]interface{}{
				"uri":      r.URL.Path,
				"method":   r.Method,
				"duration": duration,
				"status":   lw.data.status,
				"size":     lw.data.size,
			})
		})
	}
}

func AllowPostMethodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func FinishRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})
}
