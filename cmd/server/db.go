package main

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/logger"
)

const fail = 1

func mustInitDB(log logger.Logger) *sqlx.DB {
	db, err := sqlx.Connect("pgx", databaseDSN)
	if err != nil {
		log.WithError(err).Error("Unable to connect to database")
		os.Exit(fail)
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Error("Ping database is failed")
		os.Exit(fail)
	}

	return db
}
