package updates

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-metrics/internal/generated/rpc"
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

func (h *Handler) Updates(ctx context.Context, in *rpc.UpdatesIn) (*rpc.UpdatesOut, error) {
	h.log.Info("Received Updates")

	if len(in.Metrics) == 0 {
		return nil, nil
	}

	gauges := make([]storage.Gauge, 0, len(in.Metrics))
	counters := make([]storage.Counter, 0, len(in.Metrics))

	for _, m := range in.Metrics {
		switch m.Type {
		case model.TypeGauge:
			if m.Value == nil {
				h.log.Info("Invalid value")
				continue
			}

			gauges = append(gauges, storage.Gauge{
				ID:    m.Id,
				Value: *m.Value,
			})
		case model.TypeCounter:
			if m.Delta == nil {
				h.log.Info("Invalid value")
				continue
			}

			counters = append(counters, storage.Counter{
				ID:    m.Id,
				Value: *m.Delta,
			})
		}
	}

	err := h.repo.SetGauges(ctx, gauges)
	if err != nil {
		h.log.WithError(err).Error("Failed to save gauges")
		return nil, status.Error(codes.Internal, "failed to save gauges")
	}

	err = h.repo.AddCounters(ctx, counters)
	if err != nil {
		h.log.WithError(err).Error("Failed to save counters")
		return nil, status.Error(codes.Internal, "failed to save counters")
	}

	err = h.backup.Create(ctx)
	if err != nil {
		h.log.WithError(err).Error("Failed to backup data")
	}

	return &rpc.UpdatesOut{}, nil
}
