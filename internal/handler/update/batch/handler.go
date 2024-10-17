package batch

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
)

type Handler struct {
	repo   repo
	backup backup
	log    log
}

func NewHandler(repo repo, backup backup, log log) *Handler {
	return &Handler{
		repo:   repo,
		backup: backup,
		log:    log,
	}
}

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

	reqHash := r.Header.Get("HashSHA256")
	if len(reqHash) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	agentHash, err := hex.DecodeString(reqHash)
	if err != nil {
		h.log.WithError(err).Error("Error decoding hash")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	b := buf.Bytes()

	hm := hmac.New(sha256.New, []byte("secretkey"))
	hm.Write(b)
	hash := hm.Sum(nil)

	if !hmac.Equal(hash, agentHash) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var in []model.UpdateIn

	err = json.Unmarshal(b, &in)
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

	w.Header().Set("HashSHA256", hex.EncodeToString(hash))
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
