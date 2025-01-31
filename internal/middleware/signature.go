package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/securety/signature"
)

const headerHash = "HashSHA256"

// SignatureMiddleware HTTP middleware подписывает ответ.
func SignatureMiddleware(sign *signature.SignManager, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !sign.Enable() {
				next.ServeHTTP(w, r)
				return
			}

			reqSign := r.Header.Get(headerHash)
			if len(reqSign) == 0 {
				logger.Info(fmt.Sprintf("No '%s' header", headerHash))
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				logger.WithError(err).Error("Error reading request body")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()

			body := buf.Bytes()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			isValid, respSign := sign.Verify(buf.Bytes(), reqSign)
			if !isValid {
				logger.Info("Signature is not correct")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)

			w.Header().Set(headerHash, respSign)
		}

		return http.HandlerFunc(fn)
	}
}
