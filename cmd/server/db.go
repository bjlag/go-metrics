package main

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/logger"
)

func initDB(dsn string, log logger.Logger) *sqlx.DB {
	if len(dsn) == 0 {
		log.Error("dsn isn't set")
		return nil
	}

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.WithError(err).Error("Unable to connect to database")
		return nil
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = initSchema(db)
	if err != nil {
		log.WithError(err).Error("Unable to create database schema")
		return nil
	}

	log.WithField("dsn", dsn).Info("started db")

	return db
}

func initSchema(db *sqlx.DB) error {
	var schema = `
		CREATE TABLE IF NOT EXISTS gauge_metrics (
		    id varchar(100) PRIMARY KEY NOT NULL,
		    value double precision NOT NULL
		);
		
		COMMENT ON TABLE gauge_metrics IS 'Метрики типа gauge';
		COMMENT ON COLUMN gauge_metrics.id IS 'ID метрики';
		COMMENT ON COLUMN gauge_metrics.value IS 'Значение метрики';
		
		CREATE TABLE IF NOT EXISTS counter_metrics (
		    id varchar(100) PRIMARY KEY NOT NULL,
		    value int NOT NULL
		);
		
		COMMENT ON TABLE counter_metrics IS 'Метрики типа counter';
		COMMENT ON COLUMN counter_metrics.id IS 'ID метрики';
		COMMENT ON COLUMN counter_metrics.value IS 'Значение метрики';
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return err
}
