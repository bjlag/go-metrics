package middleware

import (
	"net"
	"net/http"

	"github.com/bjlag/go-metrics/internal/logger"
)

// CheckRealIPMiddleware HTTP middleware проверяет, что входящий запрос идет из разрешенной подсети.
func CheckRealIPMiddleware(trustedSubnet *net.IPNet, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet != nil {
				ip := net.ParseIP(r.Header.Get("X-Real-IP"))
				if ip == nil {
					logger.Error("Request is not contain `X-Real-IP` header. The request is rejected")
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}

				if !trustedSubnet.Contains(ip) {
					logger.WithField("IP", ip.String()).Error("Request IP is not from trusted subnet. The request is rejected")
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
