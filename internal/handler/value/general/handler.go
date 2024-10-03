package general

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Handler struct {
	repo repo
	log  log
}

func NewHandler(repo repo, log log) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Error reading request body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var in model.ValueIn

	err = json.Unmarshal(buf.Bytes(), &in)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Unmarshal error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if in.ID == "" {
		h.log.Info("Metric ID not specified")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !in.IsValid() {
		h.log.Info("Metric type is invalid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data, err := h.getResponseData(in)
	if err != nil {
		var metricNotFoundError *memory.MetricNotFoundError
		if errors.As(err, &metricNotFoundError) {
			h.log.WithField("type", in.MType).
				WithField("id", in.ID).
				Info("metric not found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		h.log.WithField("error", err.Error()).
			Error("Failed to get response data")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h Handler) getResponseData(request model.ValueIn) ([]byte, error) {
	out := &model.ValueOut{
		ID:    request.ID,
		MType: request.MType,
	}

	if request.IsGauge() {
		value, err := h.repo.GetGauge(request.ID)
		if err != nil {
			return nil, err
		}
		out.Value = &value
	}

	if request.IsCounter() {
		value, err := h.repo.GetCounter(request.ID)
		if err != nil {
			return nil, err
		}
		out.Delta = &value
	}

	data, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	return data, nil
}
