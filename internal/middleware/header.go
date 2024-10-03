package middleware

import (
	"net/http"
	"strings"
)

func SetHeaderResponse(key string, value []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, strings.Join(value, "; "))
			next.ServeHTTP(w, r)
		})
	}
}
