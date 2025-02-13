package general

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
)

// Handler обработчик HTTP запроса на обновление метрик обоих типов Counter и Gauge.
type Handler struct {
	repo   repo
	backup backup
	log    log
}

// NewHandler создает обработчик.
func NewHandler(repo repo, backup backup, log log) *Handler {
	return &Handler{
		repo:   repo,
		backup: backup,
		log:    log,
	}
}

// Handle обрабатывает HTTP запрос.
//
//	@Summary	Обновить метрику.
//	@Router		/update/ [post]
//	@Accept		json
//	@Produce	json
//	@Param		value	body		model.UpdateIn	true	"Request body"
//	@Success	200		{object}	model.UpdateOut
//	@Failure	400		"Некорректный запрос"
//	@Failure	404		"Метрика не найдена"
//	@Failure	500		"Ошибка"
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("error reading request body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	var in model.UpdateIn

	err = json.Unmarshal(buf.Bytes(), &in)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || errors.Is(err, model.ErrInvalidType) || errors.Is(err, model.ErrInvalidValue) {
			h.log.Info(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusNotFound)
			return
		}

		h.log.WithField("error", err.Error()).
			Error("Unmarshal error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.saveMetric(r.Context(), in)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to save metric")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.backup.Create(r.Context())
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to backup data")
	}

	data, err := h.getResponseData(r.Context(), in)
	if err != nil {
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

func (h *Handler) saveMetric(ctx context.Context, in model.UpdateIn) error {
	switch in.MType {
	case model.TypeCounter:
		h.repo.AddCounter(ctx, in.ID, *in.Delta)
	case model.TypeGauge:
		h.repo.SetGauge(ctx, in.ID, *in.Value)
	default:
		return fmt.Errorf("unknown metric type: %s", in.MType)
	}

	return nil
}

func (h *Handler) getResponseData(ctx context.Context, request model.UpdateIn) ([]byte, error) {
	out := &model.UpdateOut{
		ID:    request.ID,
		MType: request.MType,
	}

	if request.IsGauge() {
		value, err := h.repo.GetGauge(ctx, request.ID)
		if err != nil {
			return nil, err
		}
		out.Value = &value
	}

	if request.IsCounter() {
		value, err := h.repo.GetCounter(ctx, request.ID)
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
