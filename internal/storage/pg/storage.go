package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
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
	query := `SELECT id, value FROM gauge_metrics ORDER BY id`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		s.log.WithError(err).Error("failed to query")
		return nil
	}

	var models []modelGauge
	for rows.Next() {
		var model modelGauge
		err = rows.Scan(&model.ID, &model.Value)
		if err != nil {
			s.log.WithError(err).Error("failed to scan")
			return nil
		}

		models = append(models, model)
	}

	gauges := make(storage.Gauges, len(models))
	for _, gauge := range models {
		gauges[gauge.ID] = gauge.Value
	}

	return gauges
}

func (s Storage) GetAllCounters(ctx context.Context) storage.Counters {
	query := `SELECT id, value FROM counter_metrics ORDER BY id`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		s.log.WithError(err).Error("failed to query")
		return nil
	}

	var models []modelCounter
	for rows.Next() {
		var model modelCounter
		err = rows.Scan(&model.ID, &model.Value)
		if err != nil {
			s.log.WithError(err).Error("failed to scan")
			return nil
		}

		models = append(models, model)
	}

	counters := make(storage.Counters, len(models))
	for _, counter := range models {
		counters[counter.ID] = counter.Value
	}

	return counters
}

func (s Storage) GetGauge(ctx context.Context, id string) (float64, error) {
	query := `SELECT id, value FROM gauge_metrics WHERE id = $1`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var m modelGauge
	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&m.ID, &m.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.NewMetricNotFoundError(model.TypeGauge, id)
		}

		s.log.WithError(err).Error("failed to scan")
		return 0, err
	}

	return m.Value, nil
}

func (s Storage) SetGauge(ctx context.Context, id string, value float64) {
	query := `
		INSERT INTO gauge_metrics (id, value) VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
    		SET value = excluded.value
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, id, value)
	if err != nil {
		s.log.WithError(err).Error("error setting gauge")
		return
	}
}

func (s Storage) GetCounter(ctx context.Context, id string) (int64, error) {
	query := `SELECT id, value FROM counter_metrics WHERE id = $1`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var m modelCounter
	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&m.ID, &m.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.NewMetricNotFoundError(model.TypeCounter, id)
		}

		s.log.WithError(err).Error("failed to scan")
		return 0, err
	}

	return m.Value, nil
}

func (s Storage) AddCounter(ctx context.Context, id string, value int64) {
	query := `
		INSERT INTO counter_metrics VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
    		SET value = counter_metrics.value + $2
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("failed to prepare query")
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, id, value)
	if err != nil {
		s.log.WithError(err).Error("error setting gauge")
		return
	}
}
