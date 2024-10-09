package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type modelGauge struct {
	ID    string  `db:"id"`
	Value float64 `db:"value"`
}

type modelCounter struct {
	ID    string `db:"id"`
	Value int64  `db:"value"`
}

type Storage struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewStorage(db *sqlx.DB, log logger.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s Storage) GetAllGauges(ctx context.Context) storage.Gauges {
	var m []modelGauge
	query := `SELECT id, value FROM gauge_metrics ORDER BY id`
	err := s.db.SelectContext(ctx, &m, query)
	if err != nil {
		s.log.WithError(err).Error("error getting gauges")
		return nil
	}

	gauges := make(storage.Gauges, len(m))
	for _, gauge := range m {
		gauges[gauge.ID] = gauge.Value
	}

	return gauges
}

func (s Storage) GetAllCounters(ctx context.Context) storage.Counters {
	var m []modelCounter
	query := `SELECT id, value FROM counter_metrics ORDER BY id`
	err := s.db.SelectContext(ctx, &m, query)
	if err != nil {
		s.log.WithError(err).Error("error getting counters")
		return nil
	}

	counters := make(storage.Counters, len(m))
	for _, counter := range m {
		counters[counter.ID] = counter.Value
	}

	return counters
}

func (s Storage) GetGauge(ctx context.Context, id string) (float64, error) {
	var m modelGauge
	query := `SELECT id, value FROM gauge_metrics WHERE id = $1`
	err := s.db.SelectContext(ctx, &m, query, id)
	if err != nil {
		s.log.WithError(err).Error("error getting gauge")
		return 0, nil
	}

	return m.Value, nil
}

func (s Storage) SetGauge(ctx context.Context, id string, value float64) {
	query := `
		INSERT INTO gauge_metrics (id, value) VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
    		SET value = excluded.value
	`

	_, err := s.db.ExecContext(ctx, query, id, value)
	if err != nil {
		s.log.WithError(err).Error("error setting gauge")
		return
	}
}

func (s Storage) GetCounter(ctx context.Context, id string) (int64, error) {
	var m modelCounter
	query := `SELECT id, value FROM counter_metrics WHERE id = $1`
	err := s.db.SelectContext(ctx, &m, query, id)
	if err != nil {
		s.log.WithError(err).Error("error getting gauge")
		return 0, nil
	}

	return m.Value, nil
}

func (s Storage) AddCounter(ctx context.Context, id string, value int64) {
	query := `
		INSERT INTO counter_metrics VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
    		SET value = counter_metrics.value + $2
	`

	_, err := s.db.ExecContext(ctx, query, id, value)
	if err != nil {
		s.log.WithError(err).Error("error setting gauge")
		return
	}
}
