package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() error {
	var err error

	db, err = sql.Open("sqlite3", "./data/today.db")
	if err != nil {
		return err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return err
	}

	// Create GitHub repositories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS github_repositories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			name TEXT NOT NULL,
			avatar TEXT,
			url TEXT,
			description TEXT,
			language TEXT,
			language_color TEXT,
			stars INTEGER,
			forks INTEGER,
			current_period_stars INTEGER,
			built_by JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(author, name)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
