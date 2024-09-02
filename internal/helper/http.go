package helper

import (
	"log"
	"net/http"
	"regexp"
)

const (
	noNameMetricMsgErr      = "Metric name not specified"
	invalidMetricPathMsgErr = "Invalid metric path"
)

var (
	validRoutePattern  = regexp.MustCompile("^/update/(gauge|counter)/([a-zA-Z0-9_]+)?/(\\d+(.\\d+)?)$")
	withoutNamePattern = regexp.MustCompile("^/update/(gauge|counter)/(\\d+(.\\d+)?)$")
)

type UpdateHandler func(w http.ResponseWriter, r *http.Request, name, value string)

func MakeUpdateHandler(handler UpdateHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if withoutNamePattern.MatchString(r.URL.Path) {
			http.Error(w, noNameMetricMsgErr, http.StatusNotFound)
			return
		}

		m := validRoutePattern.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Error(w, invalidMetricPathMsgErr, http.StatusBadRequest)
			return
		}

		log.Printf("Metric received: type '%s', name '%s', value '%s'\n", m[1], m[2], m[3])

		handler(w, r, m[2], m[3])

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}
