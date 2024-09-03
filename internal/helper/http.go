package helper

import (
	"log"
	"net/http"
	"regexp"
)

const (
	invalidMetricTypeMsgErr = "Invalid metric type"
	noNameMetricMsgErr      = "Metric name not specified"
	invalidMetricPathMsgErr = "Invalid metric path"
)

var (
	validRoutePattern  = regexp.MustCompile("^/update/(gauge|counter)/([a-zA-Z0-9_]+)?/(\\d+(.\\d+)?)$")
	withoutNamePattern = regexp.MustCompile("^/update/(gauge|counter)/(\\d+(.\\d+)?)$")
)

type UpdateHandler func(w http.ResponseWriter, r *http.Request, name, value string)

type UpdateHandlerResolver struct {
	handlers map[string]UpdateHandler
}

func NewResolver() *UpdateHandlerResolver {
	return &UpdateHandlerResolver{
		handlers: make(map[string]UpdateHandler, 2),
	}
}

func (r UpdateHandlerResolver) AddHandler(name string, handler UpdateHandler) {
	r.handlers[name] = handler
}

func (r UpdateHandlerResolver) Resolve() http.HandlerFunc {
	resolver := r

	return func(w http.ResponseWriter, r *http.Request) {
		typeMetric := r.PathValue("type")
		handler, ok := resolver.handlers[typeMetric]
		if !ok {
			http.Error(w, invalidMetricTypeMsgErr, http.StatusBadRequest)
			return
		}

		nameMetric := r.PathValue("name")
		valueMetric := r.PathValue("value")

		//if r.Method != http.MethodPost {
		//	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		//	return
		//}
		//
		//if withoutNamePattern.MatchString(r.URL.Path) {
		//	http.Error(w, noNameMetricMsgErr, http.StatusNotFound)
		//	return
		//}
		//
		//m := validRoutePattern.FindStringSubmatch(r.URL.Path)
		//if m == nil {
		//	http.Error(w, invalidMetricPathMsgErr, http.StatusBadRequest)
		//	return
		//}

		log.Printf("Metric received: type '%s', nameMetric '%s', value '%s'\n", typeMetric, nameMetric, valueMetric)

		handler(w, r, nameMetric, valueMetric)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}
