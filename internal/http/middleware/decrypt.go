package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/securety/crypt"
)

func DecryptMiddleware(crypt *crypt.DecryptManager, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				logger.WithError(err).Error("Error reading body")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			decryptedBody, err := crypt.Decrypt(buf.Bytes())
			if err != nil {
				logger.WithError(err).Error("Error decrypting body")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
