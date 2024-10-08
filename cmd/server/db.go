package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/logger"
)

func initDB(log logger.Logger) *sqlx.DB {
	db, err := sqlx.Connect("pgx", databaseDSN)
	if err != nil {
		log.WithError(err).Error("Unable to connect to database")
		return nil
	}

	err = initSchema(db)
	if err != nil {
		log.WithError(err).Error("Unable to create database schema")
		return nil
	}

	return db
}

func initSchema(db *sqlx.DB) error {
	var schema = `
		CREATE TABLE IF NOT EXISTS metrics (
		    id varchar(100) PRIMARY KEY NOT NULL,
		    type varchar(50) NOT NULL,
		    delta int,
		    value double precision
		);
		
		COMMENT ON TABLE metrics IS 'Метрики приложения';
		COMMENT ON COLUMN metrics.id IS 'ID метрики';
		COMMENT ON COLUMN metrics.type IS 'Тип метрики: gauge, counter';
		COMMENT ON COLUMN metrics.delta IS 'Значение метрики counter';
		COMMENT ON COLUMN metrics.value IS 'Значение метрики gauge';
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return err
}
