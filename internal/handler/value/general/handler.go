package general

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
)

// Handler обработчик HTTP запроса на получение значения метрик типа Counter и Gauge.
type Handler struct {
	repo repo
	log  log
}

// NewHandler создает обработчик.
func NewHandler(repo repo, log log) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

// Handle обрабатывает HTTP запрос.
//
//	@Summary	Получить значение метрики.
//	@Router		/value/ [post]
//	@Accept		json
//	@Produce	json
//	@Param		value	body		model.ValueIn	true	"Request body"
//	@Success	200		{object}	model.ValueOut
//	@Failure	400		"Некорректный запрос"
//	@Failure	404		"Метрика не найдена"
//	@Failure	500		"Ошибка"
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Error reading request body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var in model.ValueIn

	err = json.Unmarshal(buf.Bytes(), &in)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || errors.Is(err, model.ErrInvalidType) {
			h.log.Info(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusNotFound)
			return
		}

		h.log.WithField("error", err.Error()).
			Error("Unmarshal error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := h.getResponseData(r.Context(), in)
	if err != nil {
		var metricNotFoundError *storage.NotFoundError
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

func (h Handler) getResponseData(ctx context.Context, in model.ValueIn) ([]byte, error) {
	out := &model.ValueOut{
		ID:    in.ID,
		MType: in.MType,
	}

	if in.IsGauge() {
		value, err := h.repo.GetGauge(ctx, in.ID)
		if err != nil {
			return nil, err
		}
		out.Value = &value
	}

	if in.IsCounter() {
		value, err := h.repo.GetCounter(ctx, in.ID)
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
