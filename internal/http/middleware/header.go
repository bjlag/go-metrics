package middleware

import (
	"net/http"
	"strings"
)

// HeaderResponseMiddleware HTTP middleware добавляет в ответ переданный заголовок.
func HeaderResponseMiddleware(key string, values ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := make([]string, len(values))
			copy(v, values)

			w.Header().Set(key, strings.Join(v, "; "))
			next.ServeHTTP(w, r)
		})
	}
}
