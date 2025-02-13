package batch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
)

// Handler обработчик HTTP запроса на обновление метрик батчами.
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
//	@Summary	Обновить набор метрик.
//	@Router		/updates/ [post]
//	@Accept		json
//	@Param		HashSHA256	header	string				false	"Подпись запроса (если включена проверка подписи)"
//	@Param		value		body	[]model.UpdateIn	true	"Request body"
//	@Success	200			"Метрики обновлены"
//	@Failure	400			"Некорректный запрос"
//	@Failure	404			"Метрика не найдена"
//	@Failure	500			"Ошибка"
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.log.WithError(err).Error("Error reading request body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	var in []model.UpdateIn

	err = json.Unmarshal(buf.Bytes(), &in)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || errors.Is(err, model.ErrInvalidType) || errors.Is(err, model.ErrInvalidValue) {
			h.log.Info(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		h.log.WithError(err).Error("Unmarshal error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.saveMetric(r.Context(), in)
	if err != nil {
		h.log.WithError(err).Error("Failed to save metric")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.backup.Create(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to backup data")
	}
}

func (h *Handler) saveMetric(ctx context.Context, in []model.UpdateIn) error {
	gauges := make([]storage.Gauge, 0, len(in))
	counters := make([]storage.Counter, 0, len(in))

	for _, u := range in {
		switch u.MType {
		case model.TypeGauge:
			if u.Value == nil {
				h.log.Info("Invalid value")
				continue
			}

			gauges = append(gauges, storage.Gauge{
				ID:    u.ID,
				Value: *u.Value,
			})
		case model.TypeCounter:
			if u.Delta == nil {
				h.log.Info("Invalid value")
				continue
			}

			counters = append(counters, storage.Counter{
				ID:    u.ID,
				Value: *u.Delta,
			})
		}
	}

	err := h.repo.SetGauges(ctx, gauges)
	if err != nil {
		h.log.WithError(err).Error("Failed to save gauges")
		return err
	}

	err = h.repo.AddCounters(ctx, counters)
	if err != nil {
		h.log.WithError(err).Error("Failed to save counters")
		return err
	}

	return nil
}
