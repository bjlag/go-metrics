package general

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bjlag/go-metrics/internal/handler/update/general/model"
)

type Handler struct {
	storage Storage
	log     Logger
}

func NewHandler(storage Storage, logger Logger) *Handler {
	return &Handler{
		storage: storage,
		log:     logger,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.Error(fmt.Sprintf("Error reading request body: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var request model.Request

	err = json.Unmarshal(buf.Bytes(), &request)
	if err != nil {
		h.log.Error(fmt.Sprintf("Unmarshal error: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if request.ID == "" {
		h.log.Info("Metric ID not specified", nil)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !request.IsValid() {
		h.log.Info("Metric type is invalid", nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.saveMetric(request)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to save metric: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response, err := h.createResponse(request)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to create response: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to write response: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h Handler) saveMetric(request model.Request) error {
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

func (h Handler) createResponse(request model.Request) ([]byte, error) {
	response := &model.Response{
		ID:    request.ID,
		MType: request.MType,
	}

	if request.IsGauge() {
		value, err := h.storage.GetGauge(request.ID)
		if err != nil {
			return nil, err
		}
		response.Value = &value
	}

	if request.IsCounter() {
		value, err := h.storage.GetCounter(request.ID)
		if err != nil {
			return nil, err
		}
		response.Delta = &value
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return data, nil
}
