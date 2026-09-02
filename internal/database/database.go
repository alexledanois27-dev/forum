package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// Open ouvre la base SQLite et crée les tables manquantes.
func Open(dbPath string) (*sql.DB, error) {
	directory := filepath.Dir(dbPath)

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	dsn := "file:" + filepath.ToSlash(dbPath) +
		"?_foreign_keys=on" +
		"&_busy_timeout=5000" +
		"&_journal_mode=WAL" +
		"&_time_format=sqlite"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("begin schema transaction: %w", err)
	}

	if _, err := tx.Exec(schemaSQL); err != nil {
		_ = tx.Rollback()
		db.Close()
		return nil, fmt.Errorf("initialize database schema: %w", err)
	}

	if err := tx.Commit(); err != nil {
		db.Close()
		return nil, fmt.Errorf("commit database schema: %w", err)
	}

	return db, nil
}
