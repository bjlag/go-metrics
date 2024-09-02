package update_common

import "net/http"

const (
	invalidMetricTypeMsgErr = "Invalid metric type"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	http.Error(w, invalidMetricTypeMsgErr, http.StatusBadRequest)
}
