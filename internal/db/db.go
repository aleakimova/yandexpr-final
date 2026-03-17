package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

var instance *sql.DB

func Get() *sql.DB {
	return instance
}

func Start(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)
	isNew := err != nil
	slog.Debug("opening sqlite database", "file", dbFile, "new", isNew)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if isNew {
		slog.Info("new database, creating scheduler table")
		if err := createTable(db); err != nil {
			db.Close()
			return nil, err
		}
	} else {
		slog.Debug("database file exists, skipping table creation")
	}

	instance = db
	slog.Info("database ready", "file", dbFile)
	return db, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY,
    date    TEXT CHECK(length(date) <= 8),
    title   TEXT CHECK(length(title) <= 100),
    comment TEXT CHECK(length(comment) <= 10000),
    repeat  TEXT CHECK(length(repeat) <= 100)
);`

func createTable(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
