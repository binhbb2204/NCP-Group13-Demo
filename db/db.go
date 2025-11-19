package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens (and creates) the SQLite database file and returns the DB handle.
func OpenDB(path string) (*sql.DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	dbPath := fmt.Sprintf("file:%s?_foreign_keys=1", path)
	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := d.Ping(); err != nil {
		d.Close()
		return nil, err
	}

	// Ensure foreign keys are enforced
	if _, err := d.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Println("warning: unable to enable foreign_keys", err)
	}

	return d, nil
}

// NeedsMigration checks if the DB has required tables (simple check for 'users').
func NeedsMigration(d *sql.DB) bool {
	var count int
	err := d.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	return err != nil || count == 0
}
