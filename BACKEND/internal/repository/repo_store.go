package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB opens a SQLite connection and runs schema.sql
func InitDB(path string) *sql.DB {
	// make sure folder exists
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	// create tables if not exist
	schema := `
	CREATE TABLE IF NOT EXISTS repos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		url TEXT NOT NULL,
		last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS collaborators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_id INTEGER,
		username TEXT NOT NULL,
		has_access BOOLEAN,
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(repo_id) REFERENCES repos(id)
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	return db
}
