package middleware

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}

	return next
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling request for %s, method %s", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
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
