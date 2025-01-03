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

// Storage обслуживает PostgreSQL хранилище.
type Storage struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewStorage создает хранилище.
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
		s.log.WithError(err).Error("Failed to prepare query")
		return nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to query")
		return nil
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		s.log.WithError(rows.Err()).Error("Failed to query")
		return nil
	}

	var models []modelGauge
	for rows.Next() {
		var m modelGauge
		err = rows.Scan(&m.ID, &m.Value)
		if err != nil {
			s.log.WithError(err).Error("Failed to scan")
			return nil
		}

		models = append(models, m)
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
		s.log.WithError(err).Error("Failed to prepare query")
		return nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to query")
		return nil
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		s.log.WithError(rows.Err()).Error("Failed to query")
		return nil
	}

	var models []modelCounter
	for rows.Next() {
		var model modelCounter
		err = rows.Scan(&model.ID, &model.Value)
		if err != nil {
			s.log.WithError(err).Error("Failed to scan")
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
		s.log.WithError(err).Error("Failed to prepare query")
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
			return 0, storage.NewMetricNotFoundError(model.TypeGauge, id, err)
		}

		s.log.WithError(err).Error("Failed to scan")
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
		s.log.WithError(err).Error("Failed to prepare query")
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, id, value)
	if err != nil {
		s.log.WithError(err).Error("Error setting gauge")
		return
	}
}

func (s Storage) SetGauges(ctx context.Context, gauges []storage.Gauge) error {
	if len(gauges) == 0 {
		return nil
	}

	models := make(map[string]modelGauge, len(gauges))
	for _, counter := range gauges {
		models[counter.ID] = modelGauge{
			ID:    counter.ID,
			Value: counter.Value,
		}
	}

	rows := make([]modelGauge, 0, len(models))
	for _, m := range models {
		rows = append(rows, m)
	}

	query := `
		INSERT INTO gauge_metrics (id, value) VALUES (:id, :value)
		ON CONFLICT (id) DO UPDATE
    		SET value = excluded.value
	`

	_, err := s.db.NamedExecContext(ctx, query, rows)
	if err != nil {
		s.log.WithError(err).Error("Error setting gauges")
		return err
	}

	return nil
}

func (s Storage) GetCounter(ctx context.Context, id string) (int64, error) {
	query := `SELECT id, value FROM counter_metrics WHERE id = $1`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		s.log.WithError(err).Error("Failed to prepare query")
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
			return 0, storage.NewMetricNotFoundError(model.TypeCounter, id, err)
		}

		s.log.WithError(err).Error("Failed to scan")
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
		s.log.WithError(err).Error("Failed to prepare query")
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, id, value)
	if err != nil {
		s.log.WithError(err).Error("Error setting gauge")
		return
	}
}

func (s Storage) AddCounters(ctx context.Context, counters []storage.Counter) error {
	if len(counters) == 0 {
		return nil
	}

	models := make(map[string]modelCounter, len(counters))
	for _, counter := range counters {
		if v, ok := models[counter.ID]; ok {
			v.Value += counter.Value
			models[counter.ID] = v
			continue
		}

		models[counter.ID] = modelCounter{
			ID:    counter.ID,
			Value: counter.Value,
		}
	}

	rows := make([]modelCounter, 0, len(models))
	for _, m := range models {
		rows = append(rows, m)
	}

	query := `
		INSERT INTO counter_metrics (id, value) VALUES (:id, :value)
		ON CONFLICT (id) DO UPDATE
    		SET value = counter_metrics.value + :value
	`

	_, err := s.db.NamedExecContext(ctx, query, rows)
	if err != nil {
		s.log.WithError(err).Error("Error setting counters")
		return err
	}

	return nil
}
