package general

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
)

type Handler struct {
	storage Storage
	backup  Backup
	log     Logger
}

func NewHandler(storage Storage, backup Backup, logger Logger) *Handler {
	return &Handler{
		storage: storage,
		backup:  backup,
		log:     logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.Error(fmt.Sprintf("Error reading request body: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	var in model.UpdateIn

	err = json.Unmarshal(buf.Bytes(), &in)
	if err != nil {
		h.log.Error(fmt.Sprintf("Unmarshal error: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if in.ID == "" {
		h.log.Info("Metric ID not specified", nil)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !in.IsValid() {
		h.log.Info("Metric type is invalid", nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.saveMetric(in)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to save metric: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.backup.Create()
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to backup data: %s", err.Error()), nil)
	}

	data, err := h.getResponseData(in)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to get response data: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to write response: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) saveMetric(request model.UpdateIn) error {
	switch request.MType {
	case model.TypeCounter:
		h.storage.AddCounter(request.ID, *request.Delta)
	case model.TypeGauge:
		h.storage.SetGauge(request.ID, *request.Value)
	default:
		return fmt.Errorf("unknown metric type: %s", request.MType)
	}

	return nil
}

func (h *Handler) getResponseData(request model.UpdateIn) ([]byte, error) {
	out := &model.UpdateOut{
		ID:    request.ID,
		MType: request.MType,
	}

	if request.IsGauge() {
		value, err := h.storage.GetGauge(request.ID)
		if err != nil {
			return nil, err
		}
		out.Value = &value
	}

	if request.IsCounter() {
		value, err := h.storage.GetCounter(request.ID)
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
