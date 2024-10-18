package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/signature"
)

const headerHash = "HashSHA256"

type Signature struct {
	sign *signature.SignManager
	log  logger.Logger
}

func NewSignature(sign *signature.SignManager, log logger.Logger) *Signature {
	return &Signature{
		sign: sign,
		log:  log,
	}
}

func (m *Signature) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.sign.Enable() {
			next.ServeHTTP(w, r)
			return
		}

		reqSign := r.Header.Get(headerHash)
		if len(reqSign) == 0 {
			m.log.Info(fmt.Sprintf("No '%s' header", headerHash))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			m.log.WithError(err).Error("Error reading request body")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_ = r.Body.Close()

		body := buf.Bytes()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		isValid, respSign := m.sign.Verify(buf.Bytes(), reqSign)
		if !isValid {
			m.log.Info("Signature is not correct")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)

		w.Header().Set(headerHash, respSign)
	})
}
