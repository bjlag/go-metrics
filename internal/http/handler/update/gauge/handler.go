package gauge

import (
	"net/http"
	"strconv"
)

// Handler обработчик HTTP запроса на обновление метрики типа Gauge.
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
//	@Summary	Обновление метрики типа Gauge.
//	@Router		/update/gauge/{name}/{value} [post]
//	@Param		name	path	string	true	"Название метрики"	example(Sys)
//	@Param		value	path	string	true	"Значение метрики"	example(1)
//	@Success	200		"Метрику обновили"
//	@Failure	400		"Некорректный запрос"
//	@Failure	404		"Метрика не найдена"
//	@Failure	500		"Ошибка"
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	nameMetric := r.PathValue("name")
	valueMetric := r.PathValue("value")

	if nameMetric == "" {
		h.log.Info("Metric name not specified")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	value, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		h.log.WithField("error", err.Error()).Error("invalid metric value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.repo.SetGauge(r.Context(), nameMetric, value)

	err = h.backup.Create(r.Context())
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("failed to backup data")
	}
}
