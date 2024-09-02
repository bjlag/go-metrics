package update_counter

import (
	"net/http"
	"strconv"
)

const (
	invalidTypeValueMsgErr = "Invalid type value of metric"
)

func Handle(w http.ResponseWriter, r *http.Request, nameMetric, valueMetric string) {
	value, err := strconv.ParseInt(valueMetric, 10, 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	_ = value
}
